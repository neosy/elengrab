package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"

	iconfig "github.com/neosy/elengrab/infrastructure/config"
	httpsrv "github.com/neosy/elengrab/internal/api/rest/server"
	httptemplates "github.com/neosy/elengrab/internal/api/rest/server/templates"
	"github.com/neosy/elengrab/internal/app/services"
	ytdlpdto "github.com/neosy/elengrab/internal/app/services/ytdlp/dto"
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
	// Default TTL for download state cache
	downloadStateCacheTTLDefault = 1 * time.Hour
	// Default TTL for channel information cache
	youtubeChannelCacheTTLDefault = 1 * 24 * time.Hour
	// Default TTL for site logo information cache
	siteLogoCacheTTLDefault = 1 * 24 * time.Hour

	// Update interval for site logo information
	logoUpdateInterval = 24 * time.Hour
	// Update interval for channel information
	channelUpdateInterval = 30 * 24 * time.Hour

	// Cache cleanup intervals
	cleanYoutubeChannelCacheInterval = 12 * time.Hour
	cleanDownloadStateCacheinterval  = 12 * time.Hour
	cleanSiteLogoCacheInterval       = 12 * time.Hour

	// defaultWorkerIdleTime is the default idle duration before a dynamic pool worker can exit.
	defaultWorkerIdleTime = 15 * time.Minute
)

func main() {
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

	cookiesDir, err := nfile.AbsPath(cfg.Elengrab.AppDir, cfg.Elengrab.CookiesDir)
	if err != nil && cfg.Elengrab.YoutubeAllowCookies {
		logger.Warn(
			"Failed to get abs path for cookies directory",
			"path", cfg.Elengrab.CookiesDir,
			"error", err,
		)
	}

	// Initialize SQLite database
	sqliteDir := absPath(cfg.Elengrab.AppDir, cfg.SQLite.DataDir)
	sqliteAuthDB, err := newDB(logger, persistence.DBAuthName.Path(sqliteDir))
	if err != nil {
		return
	}
	// Close the database on exit
	defer sqliterep.CloseDB(sqliteAuthDB)

	sqliteMainDB, err := newDB(logger, persistence.DBMainName.Path(sqliteDir))
	if err != nil {
		return
	}
	// Close the database on exit
	defer sqliterep.CloseDB(sqliteMainDB)

	sqliteMediaDB, err := newDB(logger, persistence.DBMediaName.Path(sqliteDir))
	if err != nil {
		return
	}
	// Close the database on exit
	defer sqliterep.CloseDB(sqliteMediaDB)

	// Apply all up migrations
	err = applyMigrations(logger, cfg, sqliteAuthDB, sqliteMainDB, sqliteMediaDB)
	if err != nil {
		return
	}

	// Create SQLite repositories
	dbs := map[persistence.DBName]*sql.DB{
		persistence.DBAuthName:  sqliteAuthDB,
		persistence.DBMainName:  sqliteMainDB,
		persistence.DBMediaName: sqliteMediaDB,
	}
	slRepositories := sqliterep.New(dbs)

	// Create in memory repositories
	inMemoryDeps := inmemoryrep.Dependencies{
		DownloadStateCacheTTL:  downloadStateCacheTTLDefault,
		YoutubeChannelCacheTTL: youtubeChannelCacheTTLDefault,
		SiteLogoCacheTTL:       siteLogoCacheTTLDefault,
	}
	inMemoryRepositories := inmemoryrep.New(inMemoryDeps)

	// Capture termination signals (Ctrl+C, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start worker pool
	dlManager := nworkerpool.NewDynamicWorkerPool(
		nworkerpool.WorkerPoolOptionLogger(logger),
		nworkerpool.WorkerPoolOptionMaxWorkers(cfg.Elengrab.DownloadWorkers),
		nworkerpool.WorkerPoolOptionIdleTime(defaultWorkerIdleTime),
	)
	if err := dlManager.Start(ctx); err != nil {
		logger.Error("Failed to start worker pool", "err", err)
		return
	}
	// Stop workers poll on exit
	defer dlManager.Stop()

	// Initialize services
	srvDeps := &services.Dependencies{
		DownloaderBinDir: cfg.Elengrab.DownloaderBinDir,
		DownloadsDir:     absPath(cfg.Elengrab.AppDir, cfg.Elengrab.DownloadsDir),
		YtDlpOptions: &ytdlpdto.Options{
			CookiesDir:          cookiesDir,
			YoutubeAllowCookies: cfg.Elengrab.YoutubeAllowCookies,
		},
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
			SiteLogo:       slRepositories.SiteLogo,
			User:           slRepositories.User,
			UserSession:    slRepositories.UserSession,

			// in memory
			DownloadStateCache:  inMemoryRepositories.DownloadState,
			YoutubeChannelCache: inMemoryRepositories.YoutubeChannel,
			SiteLogoCache:       inMemoryRepositories.SiteLogo,
		},
		DownloadDispetcher: dlManager,
		Services:           services,

		// Options
		AppName:               iconfig.AppName,
		DownloadsDir:          absPath(cfg.Elengrab.AppDir, cfg.Elengrab.DownloadsDir),
		DatabaseBackupsDir:    absPath(cfg.Elengrab.AppDir, cfg.SQLite.BackupsDir),
		DatabaseBackupsKeep:   cfg.Elengrab.Maintenance.DatabaseBackupsKeep,
		HistoryMode:           dtypes.MustParseHistoryMode(cfg.Elengrab.HistoryMode),
		DeleteDuplicatesScope: dtypes.MustParseUniquenessScope(cfg.Elengrab.DeleteDuplicatesScope),
		LogoUpdateInterval:    logoUpdateInterval,
		ChannelUpdateInterval: channelUpdateInterval,
	}
	uc := usecases.NewUsecases(ctx, logger, ucDeps)

	// Initialize stuck download jobs
	if err := uc.Downloader.ResetStuckJobs(ctx); err != nil {
		logger.Error("Failed to init downloader", "err", err)
		return
	}

	// Workers
	wsDeps := &workers.Dependencies{
		DownloadStateCache:  inMemoryRepositories.DownloadState,
		YoutubeChannelCache: inMemoryRepositories.YoutubeChannel,
		SiteLogoCache:       inMemoryRepositories.SiteLogo,
		// runners
		DownloaderMaintenance: uc.Downloader,
		Maintenance:           uc.Maintenance,
		// options
		IntervalUpdateHash:               cfg.Elengrab.Maintenance.IntervalUpdateHash,
		IntervalDeleteDuplicates:         cfg.Elengrab.Maintenance.IntervalDeleteDuplicates,
		IntervalDeleteMissingFiles:       cfg.Elengrab.Maintenance.IntervalDeleteMissingFiles,
		IntervalDeleteFailedDownloads:    cfg.Elengrab.Maintenance.IntervalDeleteFailedDownloads,
		EnableMoveUnmatchedFiles:         cfg.Elengrab.Maintenance.EnableMoveUnmatchedFiles,
		IntervalCleanYoutubeChannelCache: cleanYoutubeChannelCacheInterval,
		IntervalCleanDownloadStateCache:  cleanDownloadStateCacheinterval,
		IntervalCleanSiteLogoCache:       cleanSiteLogoCacheInterval,
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
