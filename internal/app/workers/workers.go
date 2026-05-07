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
	intervalUpdateHashDefault               = 8 * time.Hour
	intervalDeleteDuplicatesDefault         = 1 * time.Hour
	intervalDeleteMissingFilesDefault       = 30 * time.Minute
	intervalDeleteFailedDownloadsDefault    = 1 * time.Hour
	intervalCleanYoutubeChannelCacheDefault = 5 * time.Minute
	intervalCleanDownloadStateCacheDefault  = 20 * time.Minute
	intervalCleanSiteLogoCacheDefault       = 2 * time.Hour
	intervalCleanThumbnailCacheDefault      = 2 * time.Hour
	intervalBackupDatabaseDefault           = 1 * 24 * time.Hour
	intervalFlushWALDefault                 = 1 * time.Hour
	intervalUpdateSystemInfoDefault         = 30 * time.Minute
	intervalUpdateDBMetricsDefault          = 30 * time.Minute
)

type Dependencies struct {
	// cache in memory
	DownloadStateCache  persistence.DownloadStateCacheRepository
	YoutubeChannelCache persistence.YoutubeChannelCacheRepository
	SiteLogoCache       persistence.SiteLogoCacheRepository
	ThumbnailCache      persistence.ThumbnailCacheRepository

	// runners
	DownloaderMaintenance pworkers.DownloadMaintenanceRunner
	DBMaintenance         pworkers.DBMaintenanceRunner
	DBMMetrics            pworkers.DBMMetricsRunner
	DownloaderTask        pworkers.DownloadTaskRunner
	AuthWebMaintenance    pworkers.AuthWebMaintenanceRunner
	DownloaderMigrations  pworkers.MigrationsRunner

	// options
	MetricsEnabled                   bool
	IntervalUpdateHash               time.Duration
	IntervalDeleteDuplicates         time.Duration
	IntervalDeleteMissingFiles       time.Duration
	IntervalDeleteFailedDownloads    time.Duration
	MoveUnmatchedFilesEnabled        bool
	IntervalCleanYoutubeChannelCache time.Duration
	IntervalCleanDownloadStateCache  time.Duration
	IntervalCleanSiteLogoCache       time.Duration
	IntervalCleanThumbnailCache      time.Duration
	IntervalBackupDatabase           time.Duration
	IntervalFlushWAL                 time.Duration
	intervalUpdateSystemInfo         time.Duration
	intervalUpdateDBMetrics          time.Duration
}

func Initialize(logger *slog.Logger, ws *nworkers.Workers, deps *Dependencies) {
	now := time.Now().UTC()

	backupDatabaseStartAt := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	).Add(1 * time.Hour)

	ws.Add(nworkers.NewWorker(
		wjobs.NewStartupDatabaseJob(logger, deps.DBMaintenance),
		nworkers.WithName("StartupMaintenanceDatabase"),
		nworkers.WithMaxRuns(1),
		nworkers.WithInitialDelay(1*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewStartupAuthWebJob(logger, deps.AuthWebMaintenance),
		nworkers.WithName("StartupMaintenanceAuthWeb"),
		nworkers.WithMaxRuns(1),
		nworkers.WithInitialDelay(3*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDownloaderMigrationsJob(logger, deps.DownloaderMigrations),
		nworkers.WithName("DownloaderMigrations"),
		nworkers.WithMaxRuns(1),
		nworkers.WithInitialDelay(5*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewUpdateHashJob(logger, deps.DownloaderMaintenance),
		nworkers.WithName("UpdateHash"),
		nworkers.WithIntervalDefault(deps.IntervalUpdateHash, intervalUpdateHashDefault),
		nworkers.WithInitialDelay(7*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteDuplicatesJob(logger, deps.DownloaderMaintenance),
		nworkers.WithName("DeleteDuplicates"),
		nworkers.WithIntervalDefault(deps.IntervalDeleteDuplicates, intervalDeleteDuplicatesDefault),
		nworkers.WithInitialDelay(10*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteMissingFilesJob(logger, deps.DownloaderMaintenance, deps.MoveUnmatchedFilesEnabled),
		nworkers.WithName("DeleteMissingFiles"),
		nworkers.WithIntervalDefault(deps.IntervalDeleteMissingFiles, intervalDeleteMissingFilesDefault),
		nworkers.WithInitialDelay(15*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteFailedDownloadsJob(logger, deps.DownloaderMaintenance),
		nworkers.WithName("DeleteFailedDownloads"),
		nworkers.WithIntervalDefault(deps.IntervalDeleteFailedDownloads, intervalDeleteFailedDownloadsDefault),
		nworkers.WithInitialDelay(20*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.YoutubeChannelCache),
		nworkers.WithName("CleanYoutubeChannelCache"),
		nworkers.WithIntervalDefault(deps.IntervalCleanYoutubeChannelCache, intervalCleanYoutubeChannelCacheDefault),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.DownloadStateCache),
		nworkers.WithName("CleanDownloadStateCache"),
		nworkers.WithIntervalDefault(deps.IntervalCleanDownloadStateCache, intervalCleanDownloadStateCacheDefault),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.SiteLogoCache),
		nworkers.WithName("CleanSiteLogoCache"),
		nworkers.WithIntervalDefault(deps.IntervalCleanSiteLogoCache, intervalCleanSiteLogoCacheDefault),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.ThumbnailCache),
		nworkers.WithName("CleanThumbnailCache"),
		nworkers.WithIntervalDefault(deps.IntervalCleanThumbnailCache, intervalCleanThumbnailCacheDefault),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewbackupDatabaseJob(logger, deps.DBMaintenance),
		nworkers.WithName("DatabaseBackups"),
		nworkers.WithStartAt(backupDatabaseStartAt),
		nworkers.WithIntervalDefault(deps.IntervalBackupDatabase, intervalBackupDatabaseDefault),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewFlushWALJob(logger, deps.DBMaintenance),
		nworkers.WithName("FlusWAL"),
		nworkers.WithIntervalDefault(deps.IntervalFlushWAL, intervalFlushWALDefault),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewUpdateSystemInfoJob(logger, deps.DownloaderTask),
		nworkers.WithName("UpdateSystemInfo"),
		nworkers.WithIntervalDefault(deps.intervalUpdateSystemInfo, intervalUpdateSystemInfoDefault),
		nworkers.WithInitialDelay(1*time.Second),
	))

	if deps.MetricsEnabled {
		ws.Add(nworkers.NewWorker(
			wjobs.NewUpdateDBMetricsJob(logger, deps.DBMMetrics),
			nworkers.WithName("UpdateDBMetrics"),
			nworkers.WithIntervalDefault(deps.intervalUpdateDBMetrics, intervalUpdateDBMetricsDefault),
			nworkers.WithInitialDelay(10*time.Second),
		))
	}
}
