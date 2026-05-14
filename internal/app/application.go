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
	sqlitetypes "github.com/neosy/elengrab/internal/repository/sqlite/types"
)

const (
	// TTL for download state cache
	downloadStateCacheTTL = 20 * time.Minute
	// TTL for channel information cache
	youtubeChannelCacheTTL = 1 * 6 * time.Hour
	// TTL for site logo information cache
	siteLogoCacheTTL = 1 * 24 * time.Hour
	// TTL for thumbnail information cache
	thumbnailCacheTTL = 1 * 24 * time.Hour
	// for thumbnail file cache
	thumbnailFileCacheTTL = 60 * time.Minute

	// Update interval for site logo information
	logoUpdateInterval = 24 * time.Hour
	// Update interval for channel information
	channelUpdateInterval = 30 * 24 * time.Hour

	// Cache cleanup intervals
	cleanYoutubeChannelCacheInterval = 5 * time.Minute
	cleanDownloadStateCacheInterval  = 30 * time.Minute
	cleanSiteLogoCacheInterval       = 2 * time.Hour
	cleanThumbnailCacheInterval      = 2 * time.Hour
	cleanThumbnailFileCacheInterval  = 30 * time.Minute

	// Intervals for updating metrics
	updateSystemInfoInterval = 15 * time.Minute
	updateDBMetricsInterval  = 5 * time.Minute

	// defaultWorkerIdleTime is the default idle duration before a dynamic pool worker can exit.
	workerIdleTimeDefault = 15 * time.Minute

	// Auth session TTL and refresh interval defaults
	authSessionTTL             = 90 * 24 * time.Hour
	authSessionRefreshInterval = 10 * 24 * time.Hour

	// Directory for storing
	thumbnailsDir = "thumbnails"
	mediaDir      = "media"
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
	if err := a.initEnsureDirs(); err != nil {
		return err
	}

	cookiesDir := a.cookiesDir()

	// Initialize SQLite repositories
	slRepositories, err := a.initSQLiteRepositories()
	if err != nil {
		return err
	}

	// Initialize in memory repositories
	inMemoryRepositories := a.initInMemoryRepositories()

	// Initialize storages
	storages, err := a.initStorages()
	if err != nil {
		return err
	}
	a.DownloadsStorage = storages.Download

	// Initialize services
	srvDeps := &services.Dependencies{
		DownloaderBinDir: a.cfg.Elengrab.DownloaderBinDir,
		Storage:          storages.Download,
		YtDlpOptions: []ytdlpdto.Option{
			ytdlpdto.WithCookiesDir(cookiesDir),
			ytdlpdto.WithYoutubeAllowCookies(a.cfg.Elengrab.YoutubeAllowCookies),
		},
	}
	services, err := services.New(a.logger, srvDeps)
	if err != nil {
		a.logger.Error("Failed to initialize services", "err", err)
		return err
	}
	a.Services = services

	// Initialize worker pool
	workerPool := a.initWorkerPool()
	a.WorkerPool = workerPool

	// Initialize usecases
	ucDeps := &usecases.Dependencies{
		Repositories: usecases.DepRepositories{
			Repositories: slRepositories,

			MediaDownload:         slRepositories.MediaDownload,
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
			ThumbnailCache:      inMemoryRepositories.Thumbnail,
			ThumbnailFileCache:  inMemoryRepositories.ThumbnailFile,
		},

		Storages: usecases.DepStorages{
			Thumbnails: storages.Thumbnail,
			Downloads:  storages.Download,
		},

		DownloadDispetcher: workerPool,
		Services:           services,

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
		ThumbnailCache:      inMemoryRepositories.Thumbnail,
		ThumbnailFileCache:  inMemoryRepositories.ThumbnailFile,

		// runners
		DownloaderMaintenance: a.Usecases.Downloader,
		DBMaintenance:         a.Usecases.Maintenance,
		DBMMetrics:            slRepositories,
		DownloaderTask:        a.Usecases.Downloader,
		AuthWebStartup:        a.Usecases.AuthWeb,
		DownloaderMigrations:  a.Usecases.Downloader,

		// options
		MetricsEnabled:                 a.cfg.AdminServer.Enable && a.cfg.AdminServer.DebugConfig.EnableMetrics,
		UpdateHashInterval:             a.cfg.Elengrab.Maintenance.UpdateHashInterval,
		DeleteDuplicatesInterval:       a.cfg.Elengrab.Maintenance.DeleteDuplicatesInterval,
		DeleteMissingDownloadsInterval: a.cfg.Elengrab.Maintenance.DeleteMissingDownloadsInterval,
		DeleteFailedDownloadsInterval:  a.cfg.Elengrab.Maintenance.DeleteFailedDownloadsInterval,
		MoveUnmatchedFilesEnabled:      a.cfg.Elengrab.Maintenance.MoveUnmatchedFilesEnabled,

		CleanYoutubeChannelCacheInterval: cleanYoutubeChannelCacheInterval,
		CleanDownloadStateCacheInterval:  cleanDownloadStateCacheInterval,
		CleanSiteLogoCacheInterval:       cleanSiteLogoCacheInterval,
		CleanThumbnailCacheInterval:      cleanThumbnailCacheInterval,
		CleanThumbnailFileCacheInterval:  cleanThumbnailFileCacheInterval,

		// metrics
		UpdateSystemInfoInterval: updateSystemInfoInterval,
		UpdateDBMetricsInterval:  updateDBMetricsInterval,
	}
	a.Workers = nworkers.NewWorkers(a.logger)
	workers.Initialize(a.logger, wsDeps, a.Workers)

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

func (a *Application) initEnsureDirs() error {
	dirs := []string{
		absPath(a.cfg.Elengrab.RootDir, a.cfg.SQLite.DataDir),
		absPath(a.cfg.Elengrab.RootDir, a.cfg.SQLite.BackupsDir),
		absPath(a.cfg.Elengrab.RootDir, a.cfg.Elengrab.DownloadsDir),
		absPath(a.cfg.Elengrab.RootDir, filepath.Join(a.cfg.Elengrab.MediaDir, thumbnailsDir)),
	}
	if !ensureDirs(dirs) {
		return errors.New("one or more required directories are missing")
	}
	return nil
}

func (a *Application) cookiesDir() string {
	cookiesDir, err := nfile.AbsPath(a.cfg.Elengrab.RootDir, a.cfg.Elengrab.CookiesDir)
	if err != nil && a.cfg.Elengrab.YoutubeAllowCookies {
		a.logger.Warn(
			"Failed to get abs path for cookies directory",
			"path", a.cfg.Elengrab.CookiesDir,
			"error", err,
		)
	}
	return cookiesDir
}

func (a *Application) initSQLiteRepositories() (*sqliterep.Repositories, error) {
	var err error

	// Initialize SQLite database
	sqliteDir := absPath(a.cfg.Elengrab.RootDir, a.cfg.SQLite.DataDir)
	a.dbAuth, err = newDB(a.logger, sqliterep.AuthSchema.Path(sqliteDir))
	if err != nil {
		return nil, err
	}

	a.dbMain, err = newDB(a.logger, sqliterep.MainSchema.Path(sqliteDir))
	if err != nil {
		return nil, err
	}

	a.dbMedia, err = newDB(a.logger, sqliterep.MediaSchema.Path(sqliteDir))
	if err != nil {
		return nil, err
	}

	a.dbLink, err = newDB(a.logger, sqliterep.LinkSchema.Path(sqliteDir))
	if err != nil {
		return nil, err
	}

	var (
		authEntry  = sqlitetypes.NewDBEntry(sqliterep.AuthSchema, a.dbAuth)
		mainEntry  = sqlitetypes.NewDBEntry(sqliterep.MainSchema, a.dbMain)
		mediaEntry = sqlitetypes.NewDBEntry(sqliterep.MediaSchema, a.dbMedia)
		linkEntry  = sqlitetypes.NewDBEntry(sqliterep.LinkSchema, a.dbLink)
		entries    = []persistence.DBEntry{
			authEntry,
			mainEntry,
			mediaEntry,
			linkEntry,
		}
	)

	// Apply all up migrations
	err = applyMigrations(a.logger, a.cfg, authEntry, mainEntry, mediaEntry, linkEntry)
	if err != nil {
		return nil, err
	}

	// Create SQLite repositories
	return sqliterep.New(entries), nil
}

func (a *Application) initInMemoryRepositories() *inmemoryrep.Repositories {
	inMemoryDeps := inmemoryrep.Dependencies{
		DownloadStateCacheTTL:  downloadStateCacheTTL,
		YoutubeChannelCacheTTL: youtubeChannelCacheTTL,
		SiteLogoCacheTTL:       siteLogoCacheTTL,
		ThumbnailCacheTTL:      thumbnailCacheTTL,
		ThumbnailFileCacheTTL:  thumbnailFileCacheTTL,
	}
	return inmemoryrep.New(inMemoryDeps)
}

func (a *Application) initWorkerPool() nworkerpool.WorkerPool {
	downloadWorkers := func() uint32 {
		if a.cfg.Elengrab.DemoMode {
			return 1
		}
		return a.cfg.Elengrab.DownloadWorkers
	}
	return nworkerpool.NewDynamicWorkerPool(
		nworkerpool.WithLogger(a.logger),
		nworkerpool.WithMaxWorkers(downloadWorkers()),
		nworkerpool.WithIdleTime(workerIdleTimeDefault),
	)

}

func (a *Application) initStorages() (*fsstorage.Storages, error) {
	storages, err := fsstorage.NewStorages(
		absPath(a.cfg.Elengrab.RootDir, filepath.Join(a.cfg.Elengrab.MediaDir, thumbnailsDir)),
		absPath(a.cfg.Elengrab.RootDir, a.cfg.Elengrab.DownloadsDir),
		mediaDir,
	)
	if err != nil {
		a.logger.Error("Failed to initialize Storage", "error", err)
		return nil, err
	}
	return storages, nil
}
