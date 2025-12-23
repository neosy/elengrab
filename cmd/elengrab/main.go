package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"

	database "github.com/neosy/elengrab/db"
	iconfig "github.com/neosy/elengrab/infrastructure/config"
	httpsrv "github.com/neosy/elengrab/internal/api/rest/server"
	httptemplates "github.com/neosy/elengrab/internal/api/rest/server/templates"
	"github.com/neosy/elengrab/internal/app/services"
	"github.com/neosy/elengrab/internal/app/usecases"
	"github.com/neosy/elengrab/internal/app/workers"
	inmemoryrep "github.com/neosy/elengrab/internal/repository/in_memory"
	sqliterep "github.com/neosy/elengrab/internal/repository/sqlite"
	"github.com/neosy/elengrab/pkg/nfile"
	"github.com/neosy/elengrab/pkg/nlogger"
	"github.com/neosy/elengrab/pkg/nworkerpool"
	"github.com/neosy/elengrab/pkg/nworkers"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

const (
	databaseFileName                 = "elengrab.db"
	downloadStateCacheTTLDefault     = 1 * time.Hour
	youtubeChannelCacheTTLDefault    = 1 * 24 * time.Hour
	intervalCleanYoutubeChannelCache = 12 * time.Hour
	intervalCleanDownloadStateCache  = 12 * time.Hour
)

// absPath resolves a relative path to an absolute path using current working directory.
// Exits the program if resolving fails.
func absPath(path string) string {
	path, err := nfile.AbsPathCwd(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	return path
}

// checkDirs verifies that all given directories exist and are directories.
// Returns true if all directories are valid, false otherwise.
func checkDirs(dirs []string) bool {
	var allOk = true
	for _, dir := range dirs {
		if err := nfile.CheckDir(dir); err != nil {
			allOk = false
			log.Println(err)
		}
	}
	return allOk
}

func main() {
	var err error

	// Version output at startup
	fmt.Fprintf(os.Stderr, "%s v%s\n", iconfig.AppName, iconfig.AppVersion)

	// Load application configuration
	cfg, err := iconfig.New(uptr.String(iconfig.AppName), uptr.String(iconfig.AppVersion))
	if err != nil {
		log.Fatalln(err)
	}

	dirs := []string{
		absPath(cfg.SQLite.DataDir),
		absPath(cfg.SQLite.BackupsDir),
		absPath(cfg.Elengrab.DownloaderBinDir),
		absPath(cfg.Elengrab.AssetsDir),
		absPath(cfg.Elengrab.DownloadsDir),
	}
	if !checkDirs(dirs) {
		log.Fatal("one or more required directories are missing")
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Create a logger with Info level using HandlerOptions
	handlerOptions := &slog.HandlerOptions{
		// Set the logging level
		Level: nlogger.LevelToSlogLevel(cfg.AppConfig.LogLevel),
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))

	log.Printf("Logging level set to '%s'.\n", cfg.AppConfig.LogLevel)

	// Initialize SQLite database with migrations
	sqliteDB, err := sqliterep.InitDB(filepath.Join(absPath(cfg.SQLite.DataDir), databaseFileName))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize SQLite database: %v", err))
		return
	}
	// Close the database on exit
	defer sqliterep.CloseDB(sqliteDB)

	// Apply all up migrations
	migrations := database.NewMigrations(sqliteDB, nil)
	if err := migrations.ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize SQLite database: %v", err))
		return
	}

	// Create SQLite repositories
	slRepositories := sqliterep.New(sqliteDB)

	// Create in memory repositories
	inMemoryDeps := inmemoryrep.Dependencies{
		DownloadStateCacheTTL:  downloadStateCacheTTLDefault,
		YoutubeChannelCacheTTL: youtubeChannelCacheTTLDefault,
	}
	inMemoryRepositories := inmemoryrep.New(inMemoryDeps)

	// Capture termination signals (Ctrl+C, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start worker pool
	dlManager := nworkerpool.NewWorkerPool(logger, &nworkerpool.WorkerPoolOptions{PoolSize: cfg.Elengrab.DownloadWorkers})
	if err := dlManager.Start(ctx); err != nil {
		logger.Error("Failed to start worker pool", "err", err)
		return
	}
	// Stop workers on exit
	defer dlManager.Stop()

	// Initialize services
	srvDeps := &services.Dependencies{
		DownloaderBinDir: cfg.Elengrab.DownloaderBinDir,
		DownloadsDir:     absPath(cfg.Elengrab.DownloadsDir),
	}
	services, err := services.New(logger, srvDeps)
	if err != nil {
		logger.Error("Failed to initialize services", "err", err)
		return
	}

	// Initialize usecases
	ucDeps := &usecases.Dependencies{
		Repositories: usecases.DepRepositories{
			Database: slRepositories,

			File:           slRepositories.File,
			DownloadTask:   slRepositories.DownloadTask,
			YoutubeChannel: slRepositories.YoutubeChannel,

			// in memory
			DownloadStateCache:  inMemoryRepositories.DownloadState,
			YoutubeChannelCache: inMemoryRepositories.YoutubeChannel,
		},
		DownloadDispetcher: dlManager,
		Services:           services,

		// Options
		AppName:             iconfig.AppName,
		DownloadsDir:        absPath(cfg.Elengrab.DownloadsDir),
		DatabaseBackupsDir:  absPath(cfg.SQLite.BackupsDir),
		DatabaseBackupsKeep: cfg.Elengrab.Maintenance.DatabaseBackupsKeep,
		LoadHistory:         cfg.Elengrab.LoadHistory,
	}
	uc := usecases.NewUsecases(ctx, logger, ucDeps)

	// Init
	if err := uc.Downloader.ResetStuckJobs(ctx); err != nil {
		logger.Error("Failed to init downloader", "err", err)
		return
	}

	// Workers
	wsDeps := &workers.Dependencies{
		DownloadStateCache:  inMemoryRepositories.DownloadState,
		YoutubeChannelCache: inMemoryRepositories.YoutubeChannel,
		// runners
		DownloaderMaintenance: uc.Downloader,
		Maintenance:           uc.Maintenance,
		// options
		IntervalUpdateHash:               cfg.Elengrab.Maintenance.IntervalUpdateHash,
		IntervalDeleteDuplicates:         cfg.Elengrab.Maintenance.IntervalDeleteDuplicates,
		IntervalDeleteMissingFiles:       cfg.Elengrab.Maintenance.IntervalDeleteDuplicates,
		IntervalDeleteFailedDownloads:    cfg.Elengrab.Maintenance.IntervalDeleteFailedDownloads,
		EnableMoveUnmatchedFiles:         cfg.Elengrab.Maintenance.EnableMoveUnmatchedFiles,
		IntervalCleanYoutubeChannelCache: intervalCleanYoutubeChannelCache,
		IntervalCleanDownloadStateCache:  intervalCleanDownloadStateCache,
	}
	ws := nworkers.NewWorkers(logger, wsDeps, workers.InitWorkers)
	if err := ws.StartWorkers(ctx); err != nil {
		logger.Error("Failed to run workers", "err", err)
	}
	defer ws.StopWorkers()

	// Start FastHTTP server in a separate goroutine
	go func(ctx context.Context) {
		tmpl, err := httptemplates.LoadTemplates(absPath(cfg.Elengrab.AssetsDir))
		if err != nil {
			logger.Error(err.Error())
			cancel()
		}

		deps := &httpsrv.Dependencies{
			Usecases:     uc,
			Templates:    tmpl,
			AssetsDir:    absPath(cfg.Elengrab.AssetsDir),
			DownloadsDir: absPath(cfg.Elengrab.DownloadsDir),
		}

		httpServer := httpsrv.NewServer(logger, cfg.AppConfig.AppEnv, deps)
		err = httpServer.ListenAndServe(ctx, cfg.HTTPServer.Port)
		if err != nil {
			cancel()
		}
	}(ctx)

	// Wait for context cancellation or termination signal
	select {
	case <-ctx.Done():
		logger.ErrorContext(ctx, "Context completed, shutting down services...")
	case sig := <-sigChan:
		logger.ErrorContext(ctx, fmt.Sprintf("Termination signal received: %v, shutting down...", sig))
		cancel()
	}

	// Log final error if any
	// if err != nil {
	// 	log.Print(err)
	// }
}
