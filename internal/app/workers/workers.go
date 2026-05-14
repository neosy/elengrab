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
	updateHashIntervalDefault             = 8 * time.Hour
	deleteDuplicatesIntervalDefault       = 1 * time.Hour
	deleteMissingDownloadsIntervalDefault = 30 * time.Minute
	deleteFailedDownloadsIntervalDefault  = 1 * time.Hour

	cleanYoutubeChannelCacheIntervalDefault = 5 * time.Minute
	cleanDownloadStateCacheIntervalDefault  = 20 * time.Minute
	cleanSiteLogoCacheIntervalDefault       = 2 * time.Hour
	cleanThumbnailCacheIntervalDefault      = 2 * time.Hour
	cleanThumbnailFileCacheIntervalDefault  = 1 * time.Hour

	backupDatabaseIntervalDefault = 1 * 24 * time.Hour
	flushWALIntervalDefault       = 1 * time.Hour

	updateSystemInfoIntervalDefault = 30 * time.Minute
	updateDBMetricsIntervalDefault  = 30 * time.Minute
)

type Dependencies struct {
	// cache in memory
	DownloadStateCache  persistence.DownloadStateCacheRepository
	YoutubeChannelCache persistence.YoutubeChannelCacheRepository
	SiteLogoCache       persistence.SiteLogoCacheRepository
	ThumbnailCache      persistence.ThumbnailCacheRepository
	ThumbnailFileCache  persistence.ThumbnailFileCacheRepository

	// runners
	DownloaderMaintenance pworkers.DownloadMaintenanceRunner
	DBMaintenance         pworkers.DBMaintenanceRunner
	DBMMetrics            pworkers.DBMMetricsRunner
	DownloaderTask        pworkers.DownloadTaskRunner
	AuthWebStartup        pworkers.AuthWebStartupRunner
	DownloaderMigrations  pworkers.MigrationsRunner

	// options
	MetricsEnabled                 bool
	UpdateHashInterval             time.Duration
	DeleteDuplicatesInterval       time.Duration
	DeleteMissingDownloadsInterval time.Duration
	DeleteFailedDownloadsInterval  time.Duration
	MoveUnmatchedFilesEnabled      bool

	// caches
	CleanYoutubeChannelCacheInterval time.Duration
	CleanDownloadStateCacheInterval  time.Duration
	CleanSiteLogoCacheInterval       time.Duration
	CleanThumbnailCacheInterval      time.Duration
	CleanThumbnailFileCacheInterval  time.Duration

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
		wjobs.NewStartupDatabaseJob(logger, deps.DBMaintenance),
		nworkers.WithMaxRuns(1),
		nworkers.WithInitialDelay(1*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewStartupAuthWebJob(logger, deps.AuthWebStartup),
		nworkers.WithMaxRuns(1),
		nworkers.WithInitialDelay(3*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDownloaderMigrationsJob(logger, deps.DownloaderMigrations),
		nworkers.WithMaxRuns(1),
		nworkers.WithInitialDelay(5*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewUpdateHashJob(logger, deps.DownloaderMaintenance),
		nworkers.WithIntervalDefault(deps.UpdateHashInterval, updateHashIntervalDefault),
		nworkers.WithInitialDelay(7*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteDuplicatesJob(logger, deps.DownloaderMaintenance),
		nworkers.WithIntervalDefault(deps.DeleteDuplicatesInterval, deleteDuplicatesIntervalDefault),
		nworkers.WithInitialDelay(10*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteMissingDownloadsJob(logger, deps.DownloaderMaintenance, deps.MoveUnmatchedFilesEnabled),
		nworkers.WithIntervalDefault(deps.DeleteMissingDownloadsInterval, deleteMissingDownloadsIntervalDefault),
		nworkers.WithInitialDelay(15*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteFailedDownloadsJob(logger, deps.DownloaderMaintenance),
		nworkers.WithIntervalDefault(deps.DeleteFailedDownloadsInterval, deleteFailedDownloadsIntervalDefault),
		nworkers.WithInitialDelay(20*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.YoutubeChannelCache),
		nworkers.WithIntervalDefault(deps.CleanYoutubeChannelCacheInterval, cleanYoutubeChannelCacheIntervalDefault),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.DownloadStateCache),
		nworkers.WithIntervalDefault(deps.CleanDownloadStateCacheInterval, cleanDownloadStateCacheIntervalDefault),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.SiteLogoCache),
		nworkers.WithIntervalDefault(deps.CleanSiteLogoCacheInterval, cleanSiteLogoCacheIntervalDefault),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.ThumbnailCache),
		nworkers.WithIntervalDefault(deps.CleanThumbnailCacheInterval, cleanThumbnailCacheIntervalDefault),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.ThumbnailFileCache),
		nworkers.WithIntervalDefault(deps.CleanThumbnailFileCacheInterval, cleanThumbnailFileCacheIntervalDefault),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewbackupDatabaseJob(logger, deps.DBMaintenance),
		nworkers.WithStartAt(backupDatabaseStartAt),
		nworkers.WithIntervalDefault(deps.BackupDatabaseInterval, backupDatabaseIntervalDefault),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewFlushWALJob(logger, deps.DBMaintenance),
		nworkers.WithIntervalDefault(deps.FlushWALInterval, flushWALIntervalDefault),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewUpdateSystemInfoJob(logger, deps.DownloaderTask),
		nworkers.WithIntervalDefault(deps.UpdateSystemInfoInterval, updateSystemInfoIntervalDefault),
		nworkers.WithInitialDelay(1*time.Second),
	))

	if deps.MetricsEnabled {
		ws.Add(nworkers.NewWorker(
			wjobs.NewUpdateDBMetricsJob(logger, deps.DBMMetrics),
			nworkers.WithIntervalDefault(deps.UpdateDBMetricsInterval, updateDBMetricsIntervalDefault),
			nworkers.WithInitialDelay(10*time.Second),
		))
	}
}
