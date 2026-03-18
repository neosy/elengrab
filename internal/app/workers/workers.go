package workers

import (
	"log/slog"
	"time"

	cachejobs "github.com/neosy/elengrab/internal/app/workers/cache"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
	"github.com/neosy/elengrab/pkg/nworkers"
)

const (
	intervalUpdateHashDefault               = 8 * time.Hour
	intervalDeleteDuplicatesDefault         = 1 * time.Hour
	intervalDeleteMissingFilesDefault       = 30 * time.Minute
	intervalDeleteFailedDownloadsDefault    = 1 * time.Hour
	intervalCleanYoutubeChannelCacheDefault = 12 * time.Hour
	intervalCleanDownloadStateCacheDefault  = 12 * time.Hour
	intervalCleanSiteLogoCacheDefault       = 12 * time.Hour
	intervalBackupDatabaseDefault           = 1 * 24 * time.Hour
	intervalFlushWALDefault                 = 1 * time.Hour
)

type Dependencies struct {
	// cache in memory
	DownloadStateCache  persistence.DownloadStateCacheRepository
	YoutubeChannelCache persistence.YoutubeChannelCacheRepository
	SiteLogoCache       persistence.SiteLogoCacheRepository

	// runners
	DownloaderMaintenance pworkers.DownloadMaintenanceRunner
	Maintenance           pworkers.MaintenanceRunner

	// options
	IntervalUpdateHash               time.Duration
	IntervalDeleteDuplicates         time.Duration
	IntervalDeleteMissingFiles       time.Duration
	IntervalDeleteFailedDownloads    time.Duration
	EnableMoveUnmatchedFiles         bool
	IntervalCleanYoutubeChannelCache time.Duration
	IntervalCleanDownloadStateCache  time.Duration
	IntervalCleanSiteLogoCache       time.Duration
	IntervalBackupDatabase           time.Duration
	IntervalFlushWAL                 time.Duration
}

func InitWorkers(logger *slog.Logger, ws *nworkers.Workers, deps *Dependencies) {
	now := time.Now().UTC()

	backupDatabaseStartAt := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	).Add(1 * time.Hour)

	ws.Add(nworkers.NewWorker(
		wjobs.NewStartupDatabaseJob(logger, deps.Maintenance),
		nworkers.WorkerOptionName("StartupMaintenanceDatabase"),
		nworkers.WorkerOptionMaxRuns(1),
		nworkers.WorkerOptionOneShotDelay(1*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewUpdateHashJob(logger, deps.DownloaderMaintenance),
		nworkers.WorkerOptionName("UpdateHash"),
		nworkers.WorkerOptionIntervalWithDefault(deps.IntervalUpdateHash, intervalUpdateHashDefault),
		nworkers.WorkerOptionOneShotDelay(5*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteDuplicatesJob(logger, deps.DownloaderMaintenance),
		nworkers.WorkerOptionName("DeleteDuplicates"),
		nworkers.WorkerOptionIntervalWithDefault(deps.IntervalDeleteDuplicates, intervalDeleteDuplicatesDefault),
		nworkers.WorkerOptionOneShotDelay(10*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteMissingFilesJob(logger, deps.DownloaderMaintenance, deps.EnableMoveUnmatchedFiles),
		nworkers.WorkerOptionName("DeleteMissingFiles"),
		nworkers.WorkerOptionIntervalWithDefault(deps.IntervalDeleteMissingFiles, intervalDeleteMissingFilesDefault),
		nworkers.WorkerOptionOneShotDelay(15*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteFailedDownloadsJob(logger, deps.DownloaderMaintenance),
		nworkers.WorkerOptionName("DeleteFailedDownloads"),
		nworkers.WorkerOptionIntervalWithDefault(deps.IntervalDeleteFailedDownloads, intervalDeleteFailedDownloadsDefault),
		nworkers.WorkerOptionOneShotDelay(20*time.Second),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.YoutubeChannelCache),
		nworkers.WorkerOptionName("CleanYoutubeChannelCache"),
		nworkers.WorkerOptionIntervalWithDefault(deps.IntervalCleanYoutubeChannelCache, intervalCleanYoutubeChannelCacheDefault),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.DownloadStateCache),
		nworkers.WorkerOptionName("CleanDownloadStateCache"),
		nworkers.WorkerOptionIntervalWithDefault(deps.IntervalCleanDownloadStateCache, intervalCleanDownloadStateCacheDefault),
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.SiteLogoCache),
		nworkers.WorkerOptionName("CleanSiteLogoCache"),
		nworkers.WorkerOptionIntervalWithDefault(deps.IntervalCleanSiteLogoCache, intervalCleanSiteLogoCacheDefault),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewbackupDatabaseJob(logger, deps.Maintenance),
		nworkers.WorkerOptionName("DatabaseBackups"),
		nworkers.WorkerOptionStartAt(backupDatabaseStartAt),
		nworkers.WorkerOptionIntervalWithDefault(deps.IntervalBackupDatabase, intervalBackupDatabaseDefault),
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewFlushWALJob(logger, deps.Maintenance),
		nworkers.WorkerOptionName("FlusWAL"),
		nworkers.WorkerOptionIntervalWithDefault(deps.IntervalFlushWAL, intervalFlushWALDefault),
	))
}
