package main

import (
	"context"
	"database/sql"
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
	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	httptemplates "github.com/neosy/elengrab/internal/api/rest/server/templates"
	"github.com/neosy/elengrab/internal/app/services"
	"github.com/neosy/elengrab/internal/app/usecases"
	"github.com/neosy/elengrab/internal/app/workers"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/ports/persistence"
	inmemoryrep "github.com/neosy/elengrab/internal/repository/in_memory"
	sqliterep "github.com/neosy/elengrab/internal/repository/sqlite"
	"github.com/neosy/elengrab/pkg/nfile"
	"github.com/neosy/elengrab/pkg/nlogger"
	"github.com/neosy/elengrab/pkg/nworkerpool"
	"github.com/neosy/elengrab/pkg/nworkers"
)

const (
	downloadStateCacheTTLDefault     = 1 * time.Hour
	youtubeChannelCacheTTLDefault    = 1 * 24 * time.Hour
	intervalCleanYoutubeChannelCache = 12 * time.Hour
	intervalCleanDownloadStateCache  = 12 * time.Hour
)

// absPath resolves a relative path to an absolute path using current working directory.
// Exits the program if resolving fails.
func absPath(root, path string) string {
	path, err := nfile.AbsPath(root, path)
	if err != nil {
		log.Fatal(err.Error())
	}
	return path
}

func ensureAssets(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("cannot access directory %s: %w", path, err)
		}
	} else {
		return nfile.CheckDir(path)
	}

	if !assets.Embedded {
		return fmt.Errorf("directory does not exist: %s", path)
	}

	err = os.MkdirAll(path, 0o755)
	if err != nil {
		return fmt.Errorf("cannot create directory %s: %w", path, err)
	}

	err = assets.CopyToDir(path, assets.AssetsFS)
	if err != nil {
		return err
	}
	log.Printf("Embedded assets have been copied to %s\n", path)

	return nil
}

// ensureDirs verifies that all given directories exist and are directories.
// Returns true if all directories are valid, false otherwise.
func ensureDirs(dirs []string) bool {
	for _, dir := range dirs {
		_, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.MkdirAll(dir, 0o755)
				if err != nil {
					log.Printf("cannot create directory %s: %v", dir, err)
					return false
				}
				continue
			}
			log.Printf("cannot access directory %s: %v", dir, err)
			return false
		}
	}

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
	cfg, err := iconfig.New()
	if err != nil {
		log.Fatalln(err)
	}

	if err := ensureAssets(absPath(cfg.Elengrab.AppDir, cfg.Elengrab.AssetsDir)); err != nil {
		log.Fatalln(err)
	}

	dirs := []string{
		absPath(cfg.Elengrab.AppDir, cfg.SQLite.DataDir),
		absPath(cfg.Elengrab.AppDir, cfg.SQLite.BackupsDir),
		absPath(cfg.Elengrab.AppDir, cfg.Elengrab.DownloaderBinDir),
		absPath(cfg.Elengrab.AppDir, cfg.Elengrab.DownloadsDir),
	}
	if !ensureDirs(dirs) {
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
	sqliteMainDB, err := sqliterep.InitDB(filepath.Join(absPath(cfg.Elengrab.AppDir, cfg.SQLite.DataDir), sqliterep.DBFileName(persistence.DBMainName)))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize SQLite database: %v", persistence.DBMainName), "error", err)
		return
	}
	// Close the database on exit
	defer sqliterep.CloseDB(sqliteMainDB)
	sqliteAuthDB, err := sqliterep.InitDB(filepath.Join(absPath(cfg.Elengrab.AppDir, cfg.SQLite.DataDir), sqliterep.DBFileName(persistence.DBAuthName)))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize SQLite database: %v", persistence.DBAuthName), "error", err)
		return
	}
	// Close the database on exit
	defer sqliterep.CloseDB(sqliteAuthDB)

	// Apply all up migrations
	if err := database.NewMigrations(sqliteMainDB, persistence.DBMainName, nil).ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to migration SQLite database: %v", persistence.DBMainName), "error", err)
		return
	}
	if err := database.NewMigrations(sqliteAuthDB, persistence.DBAuthName, nil).ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to migration SQLite database: %v", persistence.DBAuthName), "error", err)
		return
	}

	// Create SQLite repositories
	dbs := map[persistence.DBName]*sql.DB{
		persistence.DBMainName: sqliteMainDB,
		persistence.DBAuthName: sqliteAuthDB,
	}
	slRepositories := sqliterep.New(dbs)

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
		DownloadsDir:     absPath(cfg.Elengrab.AppDir, cfg.Elengrab.DownloadsDir),
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
			User:           slRepositories.User,
			UserSession:    slRepositories.UserSession,

			// in memory
			DownloadStateCache:  inMemoryRepositories.DownloadState,
			YoutubeChannelCache: inMemoryRepositories.YoutubeChannel,
		},
		DownloadDispetcher: dlManager,
		Services:           services,

		// Options
		AppName:             iconfig.AppName,
		DownloadsDir:        absPath(cfg.Elengrab.AppDir, cfg.Elengrab.DownloadsDir),
		DatabaseBackupsDir:  absPath(cfg.Elengrab.AppDir, cfg.SQLite.BackupsDir),
		DatabaseBackupsKeep: cfg.Elengrab.Maintenance.DatabaseBackupsKeep,
		HistoryMode:         dtypes.MustParseHistoryMode(cfg.Elengrab.HistoryMode),
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
		tmpl, err := httptemplates.LoadTemplates(absPath(cfg.Elengrab.AppDir, cfg.Elengrab.AssetsDir))
		if err != nil {
			logger.Error(err.Error())
			cancel()
		}

		deps := &httpsrv.Dependencies{
			Usecases:     uc,
			Templates:    tmpl,
			AssetsDir:    absPath(cfg.Elengrab.AppDir, cfg.Elengrab.AssetsDir),
			DownloadsDir: absPath(cfg.Elengrab.AppDir, cfg.Elengrab.DownloadsDir),
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
