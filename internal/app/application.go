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
	ytdlpsrv "github.com/neosy/elengrab/internal/app/services/ytdlp"
	"github.com/neosy/elengrab/internal/app/usecases"
	"github.com/neosy/elengrab/internal/app/workers"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfile "github.com/neosy/elengrab/internal/pkg/filex"
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
	mediaDownloadCacheTTL = 5 * time.Minute
	// TTL for download state cache
	downloadStateCacheTTL = 10 * time.Minute
	// TTL for watch stat cache
	mediaWatchStatCacheTTL = 30 * time.Minute
	// TTL for watch position cache
	mediaUserWatchPositionCacheTTL = 30 * time.Minute
	// TTL for channel information cache
	youtubeChannelCacheTTL = 6 * time.Hour
	// TTL for site logo information cache
	siteLogoCacheTTL = 1 * 24 * time.Hour
	// TTL for thumbnail information cache
	thumbnailCacheTTL = 1 * 24 * time.Hour
	// for thumbnail file cache
	thumbnailFileCacheTTL = 60 * time.Minute
	// for static file cache
	assetFileCacheTTL = 60 * time.Hour

	// Update interval for site logo information
	logoUpdateInterval = 24 * time.Hour
	// Update interval for channel information
	channelUpdateInterval = 30 * 24 * time.Hour

	// Cache cleanup intervals
	cleanMediaDownloadCacheInterval         = 30 * time.Minute
	cleanDownloadStateCacheInterval         = 30 * time.Minute
	cleanMediaWatchStatCacheInterval        = 1 * time.Hour
	cleanMediaUserWatchPostionCacheInterval = 1 * time.Hour
	cleanYoutubeChannelCacheInterval        = 5 * time.Hour
	cleanSiteLogoCacheInterval              = 23 * time.Hour
	cleanThumbnailCacheInterval             = 23 * time.Hour
	cleanThumbnailFileCacheInterval         = 1 * time.Hour
	cleanAssetFileCacheInterval             = 1 * time.Hour

	// Intervals for updating metrics
	updateSystemInfoInterval = 15 * time.Minute
	updateDBMetricsInterval  = 5 * time.Minute

	// downloadWorkerIdleTimeDefault is the default idle duration before a dynamic pool worker can exit.
	downloadWorkerIdleTimeDefault  = 15 * time.Minute
	operationWorkerIdleTimeDefault = 5 * time.Minute

	// Auth session TTL and refresh interval defaults
	authSessionTTL             = 90 * 24 * time.Hour
	authSessionRefreshInterval = 10 * 24 * time.Hour

	// Directory for storing
	thumbnailsDir = "thumbnails"
	mediaDir      = "media"

	// Number of workers processing media watch events asynchronously.
	watchEventWorkers = 1
)

type Application struct {
	logger *slog.Logger
	cfg    *iconfig.Config

	ctx    context.Context
	cancel context.CancelFunc

	// Storages
	DownloadsStorage pstorage.DownloadsStorage

	// DB repositories
	sqLiteRepositories *sqliterep.Repositories

	// Caches
	AssetFileCacheRepository persistence.AssetFileCacheRepository

	Usecases *usecases.Usecases
	Services *services.Services

	DownloadWorkerPool   nworkerpool.WorkerPool
	OperationWorkerPool  nworkerpool.WorkerPool
	WatchEventWorkerPool nworkerpool.WorkerPool

	Workers *nworkers.Workers

	dbAuth       *sql.DB
	dbMain       *sql.DB
	dbMedia      *sql.DB
	dbLink       *sql.DB
	dbWatchEvent *sql.DB
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
	a.sqLiteRepositories = slRepositories

	// Initialize in memory repositories
	inMemoryRepositories := a.initInMemoryRepositories()
	a.AssetFileCacheRepository = inMemoryRepositories.AssetFile

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
		YtDlpOptions: []ytdlpsrv.ServiceOption{
			ytdlpsrv.WithCookiesDir(cookiesDir),
			ytdlpsrv.WithCookies(a.cfg.Elengrab.AllowCookies),
		},
	}
	services, err := services.New(a.logger, srvDeps)
	if err != nil {
		a.logger.Error("Failed to initialize services", "err", err)
		return err
	}
	a.Services = services

	// Initialize workers pool
	a.DownloadWorkerPool = a.newDownloadWorkerPool()
	a.OperationWorkerPool = a.newOperationWorkerPool()
	a.WatchEventWorkerPool = a.newWatchEventWorkerPool()

	// Initialize usecases
	ucDeps := &usecases.Dependencies{
		Repositories: usecases.DepRepositories{
			Repositories: slRepositories,

			MediaDownload:         slRepositories.MediaDownload,
			DownloadTask:          slRepositories.DownloadTask,
			DownloadDataMigration: slRepositories.DownloadDataMigration,

			MediaWatchEvent:        slRepositories.MediaWatchEvent,
			MediaUserWatchChunk:    slRepositories.MediaUserWatchChunk,
			MediaUserWatchStat:     slRepositories.MediaUserWatchStat,
			MediaWatchStat:         slRepositories.MediaWatchStat,
			MediaUserWatchPosition: slRepositories.MediaUserWatchPosition,

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
			MediaDownloadCache:          inMemoryRepositories.MediaDownload,
			DownloadStateCache:          inMemoryRepositories.DownloadState,
			MediaWatchStatCache:         inMemoryRepositories.MediaWatchStat,
			MediaUserWatchPositionCache: inMemoryRepositories.MediaUserWatchPosition,
			YoutubeChannelCache:         inMemoryRepositories.YoutubeChannel,
			SiteLogoCache:               inMemoryRepositories.SiteLogo,
			ThumbnailCache:              inMemoryRepositories.Thumbnail,
			ThumbnailFileCache:          inMemoryRepositories.ThumbnailFile,
		},

		Storages: usecases.DepStorages{
			Thumbnails: storages.Thumbnail,
			Downloads:  storages.Download,
		},

		DownloadDispetcher:   a.DownloadWorkerPool,
		OperationDispatcher:  a.OperationWorkerPool,
		WatchEventDispatcher: a.WatchEventWorkerPool,
		Services:             services,

		// Options
		AppName:  iconfig.AppName,
		AppMode:  dtypes.MustParseAppMode(a.cfg.Elengrab.Mode),
		DemoMode: a.cfg.Elengrab.DemoMode,

		AuthSessionTTL:             authSessionTTL,
		AuthSessionRefreshInterval: authSessionRefreshInterval,

		BaseURL: a.cfg.Elengrab.BaseURL,

		BaseShortURL:    strings.TrimSuffix(a.cfg.Elengrab.BaseURL, "/") + a.cfg.Elengrab.ShortLinkPrefix,
		ShortCodeLength: a.cfg.Elengrab.ShortLinkLength,
		ShortLinkTTL:    time.Duration(a.cfg.Elengrab.ShortLinkTTLDays) * 24 * time.Hour,

		DatabaseBackupsDir:  absPath(a.cfg.Elengrab.RootDir, a.cfg.SQLite.BackupsDir),
		DatabaseBackupsKeep: a.cfg.Elengrab.Maintenance.DatabaseBackupsKeep,

		DeleteDuplicatesScope: dtypes.MustParseUniquenessScope(a.cfg.Elengrab.DeleteDuplicatesScope),

		LogoUpdateInterval:    logoUpdateInterval,
		ChannelUpdateInterval: channelUpdateInterval,

		DefaultAdminLogin:    a.cfg.Elengrab.AdminLogin,
		DefaultAdminPassword: a.cfg.Elengrab.AdminPassword,
	}
	a.Usecases = usecases.NewUsecases(a.ctx, a.logger, ucDeps)

	// Start background use case workers.
	a.Usecases.Start(a.ctx)

	// Workers
	wsDeps := &workers.Dependencies{
		MediaDownloadCache:          inMemoryRepositories.MediaDownload,
		DownloadStateCache:          inMemoryRepositories.DownloadState,
		MediaWatchStatCache:         inMemoryRepositories.MediaWatchStat,
		MediaUserWatchPositionCache: inMemoryRepositories.MediaUserWatchPosition,
		YoutubeChannelCache:         inMemoryRepositories.YoutubeChannel,
		SiteLogoCache:               inMemoryRepositories.SiteLogo,
		ThumbnailCache:              inMemoryRepositories.Thumbnail,
		ThumbnailFileCache:          inMemoryRepositories.ThumbnailFile,
		AssetFileCache:              inMemoryRepositories.AssetFile,

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

		CleanMediaDownloadCacheInterval:          cleanMediaDownloadCacheInterval,
		CleanDownloadStateCacheInterval:          cleanDownloadStateCacheInterval,
		CleanMediaWatchStatCacheInterval:         cleanMediaWatchStatCacheInterval,
		CleanMediaUserWatchPositionCacheInterval: cleanMediaUserWatchPostionCacheInterval,
		CleanYoutubeChannelCacheInterval:         cleanYoutubeChannelCacheInterval,
		CleanSiteLogoCacheInterval:               cleanSiteLogoCacheInterval,
		CleanThumbnailCacheInterval:              cleanThumbnailCacheInterval,
		CleanThumbnailFileCacheInterval:          cleanThumbnailFileCacheInterval,
		CleanAssetFileCacheInterval:              cleanAssetFileCacheInterval,

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

func (a *Application) RunRequiredMigrations() error {
	a.logger.Info("Running required migrations")

	err := a.Usecases.Downloader.RunRequiredMigrations(a.ctx)
	if err != nil {
		return err
	}

	a.logger.Debug("Required migrations completed")

	return nil
}

func (a *Application) StartBackground() error {
	if err := a.DownloadWorkerPool.Start(a.ctx); err != nil {
		a.logger.Error("Failed to start download worker pool", "err", err)
		return err
	}

	if err := a.OperationWorkerPool.Start(a.ctx); err != nil {
		a.logger.Error("Failed to start operation worker pool", "err", err)
		return err
	}

	if err := a.WatchEventWorkerPool.Start(a.ctx); err != nil {
		a.logger.Error("Failed to start watch event worker pool", "err", err)
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

func (a *Application) StopWorkerPools() {
	if a.DownloadWorkerPool != nil {
		a.DownloadWorkerPool.Stop()
	}
	if a.OperationWorkerPool != nil {
		a.OperationWorkerPool.Stop()
	}
	if a.WatchEventWorkerPool != nil {
		a.WatchEventWorkerPool.Stop()
	}
}

func (a *Application) Shutdown() error {
	if a.Workers != nil {
		a.Workers.StopWorkers()
	}

	// Stop worker pools
	a.StopWorkerPools()

	// Close all database
	if a.sqLiteRepositories != nil {
		a.sqLiteRepositories.CloseAllDB()
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
	if err != nil && a.cfg.Elengrab.AllowCookies {
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

	a.dbWatchEvent, err = newDB(a.logger, sqliterep.WatchEventSchema.Path(sqliteDir))
	if err != nil {
		return nil, err
	}

	var (
		authEntry       = sqlitetypes.NewDBEntry(sqliterep.AuthSchema, a.dbAuth)
		mainEntry       = sqlitetypes.NewDBEntry(sqliterep.MainSchema, a.dbMain)
		mediaEntry      = sqlitetypes.NewDBEntry(sqliterep.MediaSchema, a.dbMedia)
		linkEntry       = sqlitetypes.NewDBEntry(sqliterep.LinkSchema, a.dbLink)
		watchEventEntry = sqlitetypes.NewDBEntry(sqliterep.WatchEventSchema, a.dbWatchEvent)

		entries = []persistence.DBEntry{
			authEntry,
			mainEntry,
			mediaEntry,
			linkEntry,
			watchEventEntry,
		}
	)

	// Apply all up migrations
	err = applyMigrations(
		a.logger, a.cfg,
		authEntry,
		mainEntry,
		mediaEntry,
		linkEntry,
		watchEventEntry,
	)
	if err != nil {
		return nil, err
	}

	// Create SQLite repositories
	return sqliterep.New(entries), nil
}

func (a *Application) initInMemoryRepositories() *inmemoryrep.Repositories {
	inMemoryDeps := inmemoryrep.Dependencies{
		MediaDownloadCacheTTL:          mediaDownloadCacheTTL,
		DownloadStateCacheTTL:          downloadStateCacheTTL,
		MediaWatchStatCacheTTL:         mediaWatchStatCacheTTL,
		MediaUserWatchPositionCacheTTL: mediaUserWatchPositionCacheTTL,
		YoutubeChannelCacheTTL:         youtubeChannelCacheTTL,
		SiteLogoCacheTTL:               siteLogoCacheTTL,
		ThumbnailCacheTTL:              thumbnailCacheTTL,
		ThumbnailFileCacheTTL:          thumbnailFileCacheTTL,
		AssetFileCacheTTL:              assetFileCacheTTL,
	}
	return inmemoryrep.New(inMemoryDeps)
}

func (a *Application) newDownloadWorkerPool() nworkerpool.WorkerPool {
	workers := func() uint32 {
		if a.cfg.Elengrab.DemoMode {
			return 1
		}
		return a.cfg.Elengrab.DownloadWorkers
	}
	return nworkerpool.NewDynamicWorkerPool(
		"Download",
		nworkerpool.WithLogger(a.logger),
		nworkerpool.WithMaxWorkers(workers()),
		nworkerpool.WithIdleTime(downloadWorkerIdleTimeDefault),
	)

}

func (a *Application) newOperationWorkerPool() nworkerpool.WorkerPool {
	return nworkerpool.NewDynamicWorkerPool(
		"Operation",
		nworkerpool.WithLogger(a.logger),
		nworkerpool.WithMaxWorkers(a.cfg.Elengrab.OperationWorkers),
		nworkerpool.WithIdleTime(operationWorkerIdleTimeDefault),
	)
}

func (a *Application) newWatchEventWorkerPool() nworkerpool.WorkerPool {
	return nworkerpool.NewWorkerPool(
		"WatchEvent",
		nworkerpool.WithLogger(a.logger),
		nworkerpool.WithMaxWorkers(watchEventWorkers),
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
