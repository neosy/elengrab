package app

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	iconfig "github.com/neosy/elengrab/internal/config"
	fsstorage "github.com/neosy/elengrab/internal/infrastructure/storage/filesystem"
	_ "modernc.org/sqlite"

	"github.com/neosy/elengrab/internal/app/services"
	ytdlpdto "github.com/neosy/elengrab/internal/app/services/ytdlp/dto"
	"github.com/neosy/elengrab/internal/app/usecases"
	"github.com/neosy/elengrab/internal/app/workers"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfile "github.com/neosy/elengrab/internal/pkg/file"
	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
	inmemoryrep "github.com/neosy/elengrab/internal/repository/in_memory"
	sqliterep "github.com/neosy/elengrab/internal/repository/sqlite"
)

const (
	// Default TTL for download state cache
	downloadStateCacheTTLDefault = 20 * time.Minute
	// Default TTL for channel information cache
	youtubeChannelCacheTTLDefault = 1 * 6 * time.Hour
	// Default TTL for site logo information cache
	siteLogoCacheTTLDefault = 1 * 24 * time.Hour

	// Update interval for site logo information
	logoUpdateInterval = 24 * time.Hour
	// Update interval for channel information
	channelUpdateInterval = 30 * 24 * time.Hour

	// Cache cleanup intervals
	cleanYoutubeChannelCacheInterval = 5 * time.Minute
	cleanDownloadStateCacheinterval  = 30 * time.Minute
	cleanSiteLogoCacheInterval       = 2 * time.Hour

	// defaultWorkerIdleTime is the default idle duration before a dynamic pool worker can exit.
	workerIdleTimeDefault = 15 * time.Minute

	// Auth session TTL and refresh interval defaults
	authSessionTTL             = 90 * 24 * time.Hour
	authSessionRefreshInterval = 10 * 24 * time.Hour

	// Directory for storing
	thumbnailsDir = "thumbnails"
)

type Application struct {
	logger *slog.Logger
	cfg    *iconfig.Config

	ctx    context.Context
	cancel context.CancelFunc

	// Storages
	DownloadsStorage pstorage.DownloadsStorage

	Usecases *usecases.Usecases
	Services *services.Services

	WorkerPool nworkerpool.WorkerPool
	Workers    *nworkers.Workers

	dbAuth  *sql.DB
	dbMain  *sql.DB
	dbMedia *sql.DB
	dbLink  *sql.DB
}

func NewApplication(logger *slog.Logger, cfg *iconfig.Config) (*Application, error) {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	app := &Application{
		logger: logger,
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	if err := app.initialize(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *Application) initialize() error {
	dirs := []string{
		absPath(a.cfg.Elengrab.RootDir, a.cfg.SQLite.DataDir),
		absPath(a.cfg.Elengrab.RootDir, a.cfg.SQLite.BackupsDir),
		absPath(a.cfg.Elengrab.RootDir, a.cfg.Elengrab.DownloadsDir),
		absPath(a.cfg.Elengrab.RootDir, filepath.Join(a.cfg.Elengrab.MediaDir, thumbnailsDir)),
	}
	if !ensureDirs(dirs) {
		return errors.New("one or more required directories are missing")
	}

	cookiesDir, err := nfile.AbsPath(a.cfg.Elengrab.RootDir, a.cfg.Elengrab.CookiesDir)
	if err != nil && a.cfg.Elengrab.YoutubeAllowCookies {
		a.logger.Warn(
			"Failed to get abs path for cookies directory",
			"path", a.cfg.Elengrab.CookiesDir,
			"error", err,
		)
	}

	// Initialize SQLite database
	sqliteDir := absPath(a.cfg.Elengrab.RootDir, a.cfg.SQLite.DataDir)
	a.dbAuth, err = newDB(a.logger, persistence.DBAuthName.Path(sqliteDir))
	if err != nil {
		return err
	}

	a.dbMain, err = newDB(a.logger, persistence.DBMainName.Path(sqliteDir))
	if err != nil {
		return err
	}

	a.dbMedia, err = newDB(a.logger, persistence.DBMediaName.Path(sqliteDir))
	if err != nil {
		return err
	}

	a.dbLink, err = newDB(a.logger, persistence.DBLinkName.Path(sqliteDir))
	if err != nil {
		return err
	}

	// Apply all up migrations
	err = applyMigrations(a.logger, a.cfg, a.dbAuth, a.dbMain, a.dbMedia, a.dbLink)
	if err != nil {
		return err
	}

	// Create SQLite repositories
	dbs := map[persistence.DBName]*sql.DB{
		persistence.DBAuthName:  a.dbAuth,
		persistence.DBMainName:  a.dbMain,
		persistence.DBMediaName: a.dbMedia,
		persistence.DBLinkName:  a.dbLink,
	}
	slRepositories := sqliterep.New(dbs)

	// Create in memory repositories
	inMemoryDeps := inmemoryrep.Dependencies{
		DownloadStateCacheTTL:  downloadStateCacheTTLDefault,
		YoutubeChannelCacheTTL: youtubeChannelCacheTTLDefault,
		SiteLogoCacheTTL:       siteLogoCacheTTLDefault,
	}
	inMemoryRepositories := inmemoryrep.New(inMemoryDeps)

	downloadWorkers := a.cfg.Elengrab.DownloadWorkers
	if a.cfg.Elengrab.DemoMode {
		downloadWorkers = 1
	}

	// Start worker pool
	a.WorkerPool = nworkerpool.NewDynamicWorkerPool(
		nworkerpool.WorkerPoolOptionLogger(a.logger),
		nworkerpool.WorkerPoolOptionMaxWorkers(downloadWorkers),
		nworkerpool.WorkerPoolOptionIdleTime(workerIdleTimeDefault),
	)

	// Initialize storages
	thumbnailsStorage, err := fsstorage.NewThumbnailsStorage(
		absPath(a.cfg.Elengrab.RootDir, filepath.Join(a.cfg.Elengrab.MediaDir, thumbnailsDir)),
	)
	if err != nil {
		a.logger.Error("Failed to initialize thumbnail Storage", "err", err)
		return err
	}
	downloadsStorage, err := fsstorage.NewDownloadsStorage(absPath(a.cfg.Elengrab.RootDir, a.cfg.Elengrab.DownloadsDir))
	if err != nil {
		a.logger.Error("Failed to initialize downloads storage", "err", err)
		return err
	}
	a.DownloadsStorage = downloadsStorage

	// Initialize services
	srvDeps := &services.Dependencies{
		DownloaderBinDir: a.cfg.Elengrab.DownloaderBinDir,
		Storage:          downloadsStorage,
		YtDlpOptions: []ytdlpdto.Option{
			ytdlpdto.WithCookiesDir(cookiesDir),
			ytdlpdto.WithYoutubeAllowCookies(a.cfg.Elengrab.YoutubeAllowCookies),
		},
	}
	a.Services, err = services.New(a.logger, srvDeps)
	if err != nil {
		a.logger.Error("Failed to initialize services", "err", err)
		return err
	}

	// Initialize usecases
	ucDeps := &usecases.Dependencies{
		Repositories: usecases.DepRepositories{
			Database: slRepositories,

			File:                  slRepositories.File,
			DownloadTask:          slRepositories.DownloadTask,
			DownloadDataMigration: slRepositories.DownloadDataMigration,

			YoutubeChannel: slRepositories.YoutubeChannel,
			SiteLogo:       slRepositories.SiteLogo,
			Thumbnail:      slRepositories.Thumbnail,

			User:        slRepositories.User,
			Role:        slRepositories.Role,
			UserRole:    slRepositories.UserRole,
			UserSession: slRepositories.UserSession,

			Link:      slRepositories.Link,
			LinkClick: slRepositories.LickClick,

			// in memory
			DownloadStateCache:  inMemoryRepositories.DownloadState,
			YoutubeChannelCache: inMemoryRepositories.YoutubeChannel,
			SiteLogoCache:       inMemoryRepositories.SiteLogo,
		},

		Storages: usecases.DepStorages{
			Thumbnails: thumbnailsStorage,
			Downloads:  downloadsStorage,
		},

		DownloadDispetcher: a.WorkerPool,
		Services:           a.Services,

		// Options
		AppName:  iconfig.AppName,
		AppMode:  dtypes.MustParseAppMode(a.cfg.Elengrab.Mode),
		DemoMode: a.cfg.Elengrab.DemoMode,

		AuthSessionTTL:             authSessionTTL,
		AuthSessionRefreshInterval: authSessionRefreshInterval,

		BaseURL: a.cfg.Elengrab.BaseURL,

		BaseShortURL:    strings.TrimSuffix(a.cfg.Elengrab.BaseURL, "/") + a.cfg.Elengrab.ShortLinkPrefix,
		ShortCodeLength: a.cfg.Elengrab.ShortLinkLength,

		DatabaseBackupsDir:  absPath(a.cfg.Elengrab.RootDir, a.cfg.SQLite.BackupsDir),
		DatabaseBackupsKeep: a.cfg.Elengrab.Maintenance.DatabaseBackupsKeep,

		DeleteDuplicatesScope: dtypes.MustParseUniquenessScope(a.cfg.Elengrab.DeleteDuplicatesScope),

		LogoUpdateInterval:    logoUpdateInterval,
		ChannelUpdateInterval: channelUpdateInterval,

		DefaultAdminLogin:    a.cfg.Elengrab.AdminLogin,
		DefaultAdminPassword: a.cfg.Elengrab.AdminPassword,
	}
	a.Usecases = usecases.NewUsecases(a.ctx, a.logger, ucDeps)

	// Workers
	wsDeps := &workers.Dependencies{
		DownloadStateCache:  inMemoryRepositories.DownloadState,
		YoutubeChannelCache: inMemoryRepositories.YoutubeChannel,
		SiteLogoCache:       inMemoryRepositories.SiteLogo,
		// runners
		DownloaderMaintenance: a.Usecases.Downloader,
		Maintenance:           a.Usecases.Maintenance,
		DownloaderTask:        a.Usecases.Downloader,
		AuthWebMaintenance:    a.Usecases.AuthWeb,
		DownloaderMigrations:  a.Usecases.Downloader,
		// options
		IntervalUpdateHash:               a.cfg.Elengrab.Maintenance.IntervalUpdateHash,
		IntervalDeleteDuplicates:         a.cfg.Elengrab.Maintenance.IntervalDeleteDuplicates,
		IntervalDeleteMissingFiles:       a.cfg.Elengrab.Maintenance.IntervalDeleteMissingFiles,
		IntervalDeleteFailedDownloads:    a.cfg.Elengrab.Maintenance.IntervalDeleteFailedDownloads,
		EnableMoveUnmatchedFiles:         a.cfg.Elengrab.Maintenance.EnableMoveUnmatchedFiles,
		IntervalCleanYoutubeChannelCache: cleanYoutubeChannelCacheInterval,
		IntervalCleanDownloadStateCache:  cleanDownloadStateCacheinterval,
		IntervalCleanSiteLogoCache:       cleanSiteLogoCacheInterval,
	}
	a.Workers = nworkers.NewWorkers(a.logger, wsDeps, workers.InitWorkers)

	return nil
}

func (a *Application) Context() context.Context {
	return a.ctx
}

func (a *Application) Cancel() {
	a.cancel()
}

func (a *Application) Logger() *slog.Logger {
	return a.logger
}

func (a *Application) StartBackground() error {
	if err := a.WorkerPool.Start(a.ctx); err != nil {
		a.logger.Error("Failed to start worker pool", "err", err)
		return err
	}

	// Initialize stuck download jobs
	if err := a.Usecases.Downloader.ResetStuckJobs(a.ctx); err != nil {
		a.logger.Error("Failed to init downloader", "err", err)
		return err
	}

	if err := a.Workers.StartWorkers(a.ctx); err != nil {
		a.logger.Error("Failed to run workers", "err", err)
		return err
	}

	return nil
}

func (a *Application) Shutdown() error {
	if a.Workers != nil {
		a.Workers.StopWorkers()
	}

	// Stop workers poll on exit
	if a.WorkerPool != nil {
		a.WorkerPool.Stop()
	}

	// Close the database on exit
	if a.dbMedia != nil {
		sqliterep.CloseDB(a.dbMedia)
	}
	if a.dbMain != nil {
		sqliterep.CloseDB(a.dbMain)
	}
	if a.dbAuth != nil {
		sqliterep.CloseDB(a.dbAuth)
	}

	return nil
}
