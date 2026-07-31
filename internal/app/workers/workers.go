package workers

import (
	"log/slog"
	"time"

	cachejobs "github.com/neosy/elengrab/internal/app/workers/cache"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
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

func Initialize(logger *slog.Logger, deps *Dependencies, ws *nworkers.Workers) {
	now := time.Now().UTC()

	backupDatabaseStartAt := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	).Add(1 * time.Hour)

	ws.Add(nworkers.NewWorker(
		wjobs.NewStartupDatabaseJob(ws.Logger(), deps.DBMaintenance),
		nworkers.WithMaxRuns(1),
		nworkers.WithInitialDelay(1*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewStartupAuthWebJob(ws.Logger(), deps.AuthWebStartup),
		nworkers.WithMaxRuns(1),
		nworkers.WithInitialDelay(3*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDownloaderMigrationsJob(ws.Logger(), deps.DownloaderMigrations),
		nworkers.WithMaxRuns(1),
		nworkers.WithInitialDelay(5*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewUpdateHashJob(ws.Logger(), deps.DownloaderMaintenance),
		nworkers.WithIntervalFallback(deps.UpdateHashInterval, defaultUpdateHashInterval),
		nworkers.WithInitialDelay(7*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteDuplicatesJob(ws.Logger(), deps.DownloaderMaintenance),
		nworkers.WithIntervalFallback(deps.DeleteDuplicatesInterval, defaultDeleteDuplicatesInterval),
		nworkers.WithInitialDelay(10*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteMissingDownloadsJob(ws.Logger(), deps.DownloaderMaintenance, deps.MoveUnmatchedFilesEnabled),
		nworkers.WithIntervalFallback(deps.DeleteMissingDownloadsInterval, defaultDeleteMissingDownloadsInterval),
		nworkers.WithInitialDelay(15*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteFailedDownloadsJob(ws.Logger(), deps.DownloaderMaintenance),
		nworkers.WithIntervalFallback(deps.DeleteFailedDownloadsInterval, defaultDeleteFailedDownloadsInterval),
		nworkers.WithInitialDelay(20*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.YoutubeChannelCache),
		nworkers.WithIntervalFallback(deps.CleanYoutubeChannelCacheInterval, defaultCleanYoutubeChannelCacheInterval),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.MediaDownloadCache),
		nworkers.WithIntervalFallback(deps.CleanMediaDownloadCacheInterval, defaultCleanMediaDownloadCacheInterval),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.DownloadStateCache),
		nworkers.WithIntervalFallback(deps.CleanDownloadStateCacheInterval, defaultCleanDownloadStateCacheInterval),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.MediaWatchStatCache),
		nworkers.WithIntervalFallback(deps.CleanMediaWatchStatCacheInterval, defaultCleanMediaWatchStatCacheInterval),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.MediaUserWatchPositionCache),
		nworkers.WithIntervalFallback(deps.CleanMediaUserWatchPositionCacheInterval, defaultCleanMediaUserWatchPositionCacheInterval),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.SiteLogoCache),
		nworkers.WithIntervalFallback(deps.CleanSiteLogoCacheInterval, defaultCleanSiteLogoCacheInterval),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.ThumbnailCache),
		nworkers.WithIntervalFallback(deps.CleanThumbnailCacheInterval, defaultCleanThumbnailCacheInterval),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.ThumbnailFileCache),
		nworkers.WithIntervalFallback(deps.CleanThumbnailFileCacheInterval, defaultCleanThumbnailFileCacheInterval),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(ws.Logger(), deps.AssetFileCache),
		nworkers.WithIntervalFallback(deps.CleanAssetFileCacheInterval, defaultCleanAssetFileCacheInterval),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewbackupDatabaseJob(ws.Logger(), deps.DBMaintenance),
		nworkers.WithStartAt(backupDatabaseStartAt),
		nworkers.WithIntervalFallback(deps.BackupDatabaseInterval, defaultBackupDatabaseInterval),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewFlushWALJob(ws.Logger(), deps.DBMaintenance),
		nworkers.WithIntervalFallback(deps.FlushWALInterval, defaultFlushWALInterval),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewUpdateSystemInfoJob(ws.Logger(), deps.DownloaderTask),
		nworkers.WithIntervalFallback(deps.UpdateSystemInfoInterval, defaultUpdateSystemInfoInterval),
		nworkers.WithInitialDelay(1*time.Second),
	))

	if deps.MetricsEnabled {
		ws.Add(nworkers.NewWorker(
			wjobs.NewUpdateDBMetricsJob(ws.Logger(), deps.DBMMetrics),
			nworkers.WithIntervalFallback(deps.UpdateDBMetricsInterval, defaultUpdateDBMetricsInterval),
			nworkers.WithInitialDelay(10*time.Second),
		))
	}
}
