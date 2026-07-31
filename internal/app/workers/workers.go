package appworkers

import (
	"log/slog"
	"time"

	cachejobs "github.com/neosy/elengrab/internal/app/workers/cache"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	"github.com/neosy/elengrab/internal/pkg/workers"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

const (
	defaultUpdateHashInterval             = 8 * time.Hour
	defaultDeleteDuplicatesInterval       = 1 * time.Hour
	defaultDeleteMissingDownloadsInterval = 30 * time.Minute
	defaultDeleteFailedDownloadsInterval  = 1 * time.Hour

	defaultCleanMediaDownloadCacheInterval          = 30 * time.Minute
	defaultCleanDownloadStateCacheInterval          = 20 * time.Minute
	defaultCleanMediaWatchStatCacheInterval         = 1 * time.Hour
	defaultCleanMediaUserWatchPositionCacheInterval = 1 * time.Hour
	defaultCleanYoutubeChannelCacheInterval         = 6 * time.Hour
	defaultCleanSiteLogoCacheInterval               = 24 * time.Hour
	defaultCleanThumbnailCacheInterval              = 24 * time.Hour
	defaultCleanThumbnailFileCacheInterval          = 1 * time.Hour
	defaultCleanAssetFileCacheInterval              = 1 * time.Hour

	defaultBackupDatabaseInterval = 1 * 24 * time.Hour
	defaultFlushWALInterval       = 1 * time.Hour

	defaultUpdateSystemInfoInterval = 30 * time.Minute
	defaultUpdateDBMetricsInterval  = 30 * time.Minute
)

type Dependencies struct {
	// cache in memory
	MediaDownloadCache          persistence.MediaDownloadCacheRepository
	DownloadStateCache          persistence.DownloadStateCacheRepository
	MediaWatchStatCache         persistence.MediaWatchStatCacheRepository
	MediaUserWatchPositionCache persistence.MediaUserWatchPositionCacheRepository
	YoutubeChannelCache         persistence.YoutubeChannelCacheRepository
	SiteLogoCache               persistence.SiteLogoCacheRepository
	ThumbnailCache              persistence.ThumbnailCacheRepository
	ThumbnailFileCache          persistence.ThumbnailFileCacheRepository
	AssetFileCache              persistence.AssetFileCacheRepository

	// runners
	DownloaderMaintenance pworkers.DownloadMaintenanceRunner
	DBMaintenance         pworkers.DBMaintenanceRunner
	DBMMetrics            pworkers.DBMMetricsRunner
	DownloaderTask        pworkers.DownloadTaskRunner
	AuthWebStartup        pworkers.AuthWebStartupRunner
	DownloaderMigrations  pworkers.MigrationsRunner

	// options
	MetricsEnabled                 bool
	MoveUnmatchedFilesEnabled      bool
	UpdateHashInterval             time.Duration
	DeleteDuplicatesInterval       time.Duration
	DeleteMissingDownloadsInterval time.Duration
	DeleteFailedDownloadsInterval  time.Duration

	// caches
	CleanYoutubeChannelCacheInterval         time.Duration
	CleanMediaDownloadCacheInterval          time.Duration
	CleanDownloadStateCacheInterval          time.Duration
	CleanMediaWatchStatCacheInterval         time.Duration
	CleanMediaUserWatchPositionCacheInterval time.Duration
	CleanSiteLogoCacheInterval               time.Duration
	CleanThumbnailCacheInterval              time.Duration
	CleanThumbnailFileCacheInterval          time.Duration
	CleanAssetFileCacheInterval              time.Duration

	// db
	BackupDatabaseInterval time.Duration
	FlushWALInterval       time.Duration

	// metrics
	UpdateSystemInfoInterval time.Duration
	UpdateDBMetricsInterval  time.Duration
}

func Initialize(logger *slog.Logger, deps *Dependencies, ws *workers.Workers) {
	now := time.Now().UTC()

	backupDatabaseStartAt := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	).Add(1 * time.Hour)

	ws.Add(workers.NewWorker(
		wjobs.NewStartupDatabaseJob(ws.Logger(), deps.DBMaintenance),
		workers.WithMaxRuns(1),
		workers.WithInitialDelay(1*time.Second),
	))

	ws.Add(workers.NewWorker(
		wjobs.NewStartupAuthWebJob(ws.Logger(), deps.AuthWebStartup),
		workers.WithMaxRuns(1),
		workers.WithInitialDelay(3*time.Second),
	))

	ws.Add(workers.NewWorker(
		wjobs.NewDownloaderMigrationsJob(ws.Logger(), deps.DownloaderMigrations),
		workers.WithMaxRuns(1),
		workers.WithInitialDelay(5*time.Second),
	))

	ws.Add(workers.NewWorker(
		wjobs.NewUpdateHashJob(ws.Logger(), deps.DownloaderMaintenance),
		workers.WithIntervalFallback(deps.UpdateHashInterval, defaultUpdateHashInterval),
		workers.WithInitialDelay(7*time.Second),
	))

	ws.Add(workers.NewWorker(
		wjobs.NewDeleteDuplicatesJob(ws.Logger(), deps.DownloaderMaintenance),
		workers.WithIntervalFallback(deps.DeleteDuplicatesInterval, defaultDeleteDuplicatesInterval),
		workers.WithInitialDelay(10*time.Second),
	))

	ws.Add(workers.NewWorker(
		wjobs.NewDeleteMissingDownloadsJob(ws.Logger(), deps.DownloaderMaintenance, deps.MoveUnmatchedFilesEnabled),
		workers.WithIntervalFallback(deps.DeleteMissingDownloadsInterval, defaultDeleteMissingDownloadsInterval),
		workers.WithInitialDelay(15*time.Second),
	))

	ws.Add(workers.NewWorker(
		wjobs.NewDeleteFailedDownloadsJob(ws.Logger(), deps.DownloaderMaintenance),
		workers.WithIntervalFallback(deps.DeleteFailedDownloadsInterval, defaultDeleteFailedDownloadsInterval),
		workers.WithInitialDelay(20*time.Second),
	))

	ws.Add(workers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.YoutubeChannelCache),
		workers.WithIntervalFallback(deps.CleanYoutubeChannelCacheInterval, defaultCleanYoutubeChannelCacheInterval),
	))

	ws.Add(workers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.MediaDownloadCache),
		workers.WithIntervalFallback(deps.CleanMediaDownloadCacheInterval, defaultCleanMediaDownloadCacheInterval),
	))

	ws.Add(workers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.DownloadStateCache),
		workers.WithIntervalFallback(deps.CleanDownloadStateCacheInterval, defaultCleanDownloadStateCacheInterval),
	))

	ws.Add(workers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.MediaWatchStatCache),
		workers.WithIntervalFallback(deps.CleanMediaWatchStatCacheInterval, defaultCleanMediaWatchStatCacheInterval),
	))

	ws.Add(workers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.MediaUserWatchPositionCache),
		workers.WithIntervalFallback(deps.CleanMediaUserWatchPositionCacheInterval, defaultCleanMediaUserWatchPositionCacheInterval),
	))

	ws.Add(workers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.SiteLogoCache),
		workers.WithIntervalFallback(deps.CleanSiteLogoCacheInterval, defaultCleanSiteLogoCacheInterval),
	))

	ws.Add(workers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.ThumbnailCache),
		workers.WithIntervalFallback(deps.CleanThumbnailCacheInterval, defaultCleanThumbnailCacheInterval),
	))

	ws.Add(workers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.ThumbnailFileCache),
		workers.WithIntervalFallback(deps.CleanThumbnailFileCacheInterval, defaultCleanThumbnailFileCacheInterval),
	))

	ws.Add(workers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.AssetFileCache),
		workers.WithIntervalFallback(deps.CleanAssetFileCacheInterval, defaultCleanAssetFileCacheInterval),
	))

	ws.Add(workers.NewWorker(
		wjobs.NewbackupDatabaseJob(ws.Logger(), deps.DBMaintenance),
		workers.WithStartAt(backupDatabaseStartAt),
		workers.WithIntervalFallback(deps.BackupDatabaseInterval, defaultBackupDatabaseInterval),
	))

	ws.Add(workers.NewWorker(
		wjobs.NewFlushWALJob(ws.Logger(), deps.DBMaintenance),
		workers.WithIntervalFallback(deps.FlushWALInterval, defaultFlushWALInterval),
	))

	ws.Add(workers.NewWorker(
		wjobs.NewUpdateSystemInfoJob(ws.Logger(), deps.DownloaderTask),
		workers.WithIntervalFallback(deps.UpdateSystemInfoInterval, defaultUpdateSystemInfoInterval),
		workers.WithInitialDelay(1*time.Second),
	))

	if deps.MetricsEnabled {
		ws.Add(workers.NewWorker(
			wjobs.NewUpdateDBMetricsJob(ws.Logger(), deps.DBMMetrics),
			workers.WithIntervalFallback(deps.UpdateDBMetricsInterval, defaultUpdateDBMetricsInterval),
			workers.WithInitialDelay(10*time.Second),
		))
	}
}
