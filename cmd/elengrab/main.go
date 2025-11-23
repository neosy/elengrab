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

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"

	database "github.com/neosy/elengrab/db"
	iconfig "github.com/neosy/elengrab/infrastructure/config"
	httpsrv "github.com/neosy/elengrab/internal/api/rest/server"
	"github.com/neosy/elengrab/internal/app/usecases"
	"github.com/neosy/elengrab/internal/app/workers"
	inmemoryrep "github.com/neosy/elengrab/internal/repository/in_memory"
	sqliterep "github.com/neosy/elengrab/internal/repository/sqlite"
	"github.com/neosy/elengrab/internal/services"
	"github.com/neosy/elengrab/pkg/nlogger"
	"github.com/neosy/elengrab/pkg/nworkers"
	"github.com/neosy/elengrab/pkg/workerpool"
)

func main() {
	var err error

	// Load application configuration
	cfg := iconfig.New()

	fmt.Printf("%s v%s\n", cfg.AppName, cfg.AppConfig.Version)

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
	sqliteDB, err := sqliterep.InitDB(filepath.Join(cfg.SQLite.DataDir, "elengrab.db"))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize SQLite database: %v", err))
		return
	}
	// Close the database on exit
	defer sqliteDB.Close()

	// Apply all up migrations
	migrations := database.NewMigrations(sqliteDB, nil)
	if err := migrations.ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize SQLite database: %v", err))
		return
	}

	// Create SQLite repositories
	slRepositories := sqliterep.New(sqliteDB)

	// Create SQLite repositories
	inMemoryRepositories := inmemoryrep.New()

	// Capture termination signals (Ctrl+C, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start workers
	dlManager := workerpool.NewManager(logger, &workerpool.ManagerOptions{WorkerCount: 3})
	if err := dlManager.Start(ctx); err != nil {
		logger.Error("Failed to start worker pool", "err", err)
		return
	}
	// Stop workers on exit
	defer dlManager.Stop()

	// Initialize services
	srvDeps := &services.Dependencies{
		BinDir:       cfg.Elengrab.BinDir,
		DownloadsDir: cfg.Elengrab.DownloadsDir,
	}
	services, err := services.New(logger, srvDeps)
	if err != nil {
		logger.Error("Failed to initialize services", "err", err)
		return
	}

	// Initialize usecases
	ucDeps := &usecases.Dependencies{
		Repositories: usecases.DepRepositories{
			File:          slRepositories.File,
			DownloadTask:  slRepositories.DownloadTask,
			DownloadState: inMemoryRepositories.DownloadState,
		},
		DownloadDispetcher: dlManager,
		Services:           services,

		// Options
		DownloadsDir: cfg.Elengrab.DownloadsDir,
		LoadHistory:  cfg.Elengrab.LoadHistory,
	}
	uc := usecases.NewUsecases(logger, ucDeps)

	// Init
	if err := uc.Downloader.ResetStuckJobs(ctx); err != nil {
		logger.Error("Failed to init downloader", "err", err)
		return
	}

	// Workers
	wsDeps := &workers.Dependencies{
		Downloader: uc.Downloader,
		// options
		IntervalUpdateHash:       cfg.Elengrab.Maintenance.IntervalUpdateHash,
		IntervalDeleteDuplicates: cfg.Elengrab.Maintenance.IntervalDeleteDuplicates,
	}
	ws := nworkers.NewWorkers(logger, wsDeps, workers.InitWorkers)
	if err := ws.StartWorkers(ctx); err != nil {
		logger.Error("Failed to run workers", "err", err)
	}
	defer ws.StopWorkers()

	// Start FastHTTP server in a separate goroutine
	go func(ctx context.Context) {
		deps := &httpsrv.Dependencies{
			Usecases:  uc,
			AssetsDir: cfg.Elengrab.AssetsDir,
		}

		httpServer := httpsrv.NewServer(logger, deps)
		err = httpServer.ListenAndServe(ctx, cfg.HTMXServer.Port)
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
	if err != nil {
		log.Print(err)
	}
}
