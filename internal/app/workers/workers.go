package workers

import (
	"log/slog"
	"time"

	cachejobs "github.com/neosy/elengrab/internal/app/workers/cache"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
	"github.com/neosy/elengrab/pkg/nworkers"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

const (
	intervalUpdateHashDefault               = 8 * time.Hour
	intervalDeleteDuplicatesDefault         = 1 * time.Hour
	intervalDeleteMissingFilesDefault       = 30 * time.Minute
	intervalDeleteFailedDownloadsDefault    = 1 * time.Hour
	intervalCleanYoutubeChannelCacheDefault = 12 * time.Hour
	intervalCleanDownloadStateCacheDefault  = 12 * time.Hour
	intervalBackupDatabaseDefault           = 1 * 24 * time.Hour
	intervalFlushWALDefault                 = 1 * time.Hour
)

type Dependencies struct {
	// cache in memory
	DownloadStateCache  persistence.DownloadStateCacheRepository
	YoutubeChannelCache persistence.YoutubeChannelCacheRepository

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
	IntervalBackupDatabase           time.Duration
	IntervalFlushWAL                 time.Duration
}

func InitWorkers(logger *slog.Logger, ws *nworkers.Workers, deps *Dependencies) {
	var (
		intervalUpdateHash               = nworkers.NewInterval(intervalUpdateHashDefault, deps.IntervalUpdateHash)
		intervalDeleteDuplicates         = nworkers.NewInterval(intervalDeleteDuplicatesDefault, deps.IntervalDeleteDuplicates)
		intervalDeleteMissingFiles       = nworkers.NewInterval(intervalDeleteMissingFilesDefault, deps.IntervalDeleteMissingFiles)
		intervalDeleteFailedDownloads    = nworkers.NewInterval(intervalDeleteFailedDownloadsDefault, deps.IntervalDeleteFailedDownloads)
		intervalCleanYoutubeChannelCache = nworkers.NewInterval(intervalCleanYoutubeChannelCacheDefault, deps.IntervalCleanYoutubeChannelCache)
		intervalCleanDownloadStateCache  = nworkers.NewInterval(intervalCleanDownloadStateCacheDefault, deps.IntervalCleanDownloadStateCache)
		intervalBackupDatabase           = nworkers.NewInterval(intervalBackupDatabaseDefault, deps.IntervalBackupDatabase)
		intervalFlushWAL                 = nworkers.NewInterval(intervalFlushWALDefault, deps.IntervalFlushWAL)
	)

	now := time.Now()
	backupDatabaseStartAt := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	).Add(1 * time.Hour)

	ws.Add(nworkers.NewWorker(
		wjobs.NewUpdateHashJob(logger, deps.DownloaderMaintenance),
		&nworkers.WorkerOptions{
			Name:         "UpdateHash",
			Interval:     intervalUpdateHash.DurationPtr(),
			OneShotDelay: uptr.Any(3 * time.Second),
		},
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteDuplicatesJob(logger, deps.DownloaderMaintenance),
		&nworkers.WorkerOptions{
			Name:         "DeleteDuplicates",
			Interval:     intervalDeleteDuplicates.DurationPtr(),
			OneShotDelay: uptr.Any(10 * time.Second),
		},
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteMissingFilesJob(logger, deps.DownloaderMaintenance, deps.EnableMoveUnmatchedFiles),
		&nworkers.WorkerOptions{
			Name:         "DeleteMissingFiles",
			Interval:     intervalDeleteMissingFiles.DurationPtr(),
			OneShotDelay: uptr.Any(20 * time.Second),
		},
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteFailedDownloadsJob(logger, deps.DownloaderMaintenance),
		&nworkers.WorkerOptions{
			Name:         "DeleteFailedDownloads",
			Interval:     intervalDeleteFailedDownloads.DurationPtr(),
			OneShotDelay: uptr.Any(30 * time.Second),
		},
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.YoutubeChannelCache),
		&nworkers.WorkerOptions{
			Name:     "CleanYoutubeChannelCache",
			Interval: intervalCleanYoutubeChannelCache.DurationPtr(),
		},
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(logger, deps.DownloadStateCache),
		&nworkers.WorkerOptions{
			Name:     "CleanDownloadStateCache",
			Interval: intervalCleanDownloadStateCache.DurationPtr(),
		},
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewbackupDatabaseJob(logger, deps.Maintenance),
		&nworkers.WorkerOptions{
			Name:     "DatabaseBackups",
			StartAt:  uptr.Any(backupDatabaseStartAt),
			Interval: intervalBackupDatabase.DurationPtr(),
		},
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewFlushWALJob(logger, deps.Maintenance),
		&nworkers.WorkerOptions{
			Name:     "FlusWAL",
			Interval: intervalFlushWAL.DurationPtr(),
		},
	))
}
