package workers

import (
	"time"

	ytdownloader "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader"
	cachejobs "github.com/neosy/elengrab/internal/app/workers/cache"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/pkg/nworkers"
)

const (
	intervalUpdateHashDefault               = 8 * time.Hour
	intervalDeleteDuplicatesDefault         = 1 * time.Hour
	intervalDeleteMissingFilesDefault       = 30 * time.Minute
	intervalDeleteFailedDownloadsDefault    = 1 * time.Hour
	intervalCleanYoutubeChannelCacheDefault = 1 * time.Hour
	intervalCleanDownloadStateCacheDefault  = 12 * time.Hour
)

type Dependencies struct {
	// cache in memory
	DownloadStateCache  persistence.DownloadStateCacheRepository
	YoutubeChannelCache persistence.YoutubeChannelCacheRepository

	// usecases
	Downloader *ytdownloader.YouTubeDownloader

	// options
	IntervalUpdateHash               time.Duration
	IntervalDeleteDuplicates         time.Duration
	IntervalDeleteMissingFiles       time.Duration
	IntervalDeleteFailedDownloads    time.Duration
	EnableMoveUnmatchedFiles         bool
	IntervalCleanYoutubeChannelCache time.Duration
	IntervalCleanDownloadStateCache  time.Duration
}

func InitWorkers(ws *nworkers.Workers, deps *Dependencies) {
	var (
		intervalUpdateHash               = nworkers.NewInterval(intervalUpdateHashDefault, deps.IntervalUpdateHash)
		intervalDeleteDuplicates         = nworkers.NewInterval(intervalDeleteDuplicatesDefault, deps.IntervalDeleteDuplicates)
		intervalDeleteMissingFiles       = nworkers.NewInterval(intervalDeleteMissingFilesDefault, deps.IntervalDeleteMissingFiles)
		intervalDeleteFailedDownloads    = nworkers.NewInterval(intervalDeleteFailedDownloadsDefault, deps.IntervalDeleteFailedDownloads)
		intervalCleanYoutubeChannelCache = nworkers.NewInterval(intervalCleanYoutubeChannelCacheDefault, deps.IntervalCleanYoutubeChannelCache)
		intervalCleanDownloadStateCache  = nworkers.NewInterval(intervalCleanDownloadStateCacheDefault, deps.IntervalCleanDownloadStateCache)
	)

	ws.Add(nworkers.NewWorker(
		wjobs.NewUpdateHashJob(deps.Downloader),
		&nworkers.WorkerOptions{
			Name:       "UpdateHash",
			Interval:   intervalUpdateHash.DurationPtr(),
			FirstDelay: 3 * time.Second,
		},
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteDuplicatesJob(deps.Downloader),
		&nworkers.WorkerOptions{
			Name:       "DeleteDuplicates",
			Interval:   intervalDeleteDuplicates.DurationPtr(),
			FirstDelay: 10 * time.Second,
		},
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteMissingFilesJob(deps.Downloader, deps.EnableMoveUnmatchedFiles),
		&nworkers.WorkerOptions{
			Name:       "DeleteMissingFiles",
			Interval:   intervalDeleteMissingFiles.DurationPtr(),
			FirstDelay: 20 * time.Second,
		},
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteFailedDownloadsJob(deps.Downloader),
		&nworkers.WorkerOptions{
			Name:       "DeleteFailedDownloads",
			Interval:   intervalDeleteFailedDownloads.DurationPtr(),
			FirstDelay: 30 * time.Second,
		},
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(deps.YoutubeChannelCache),
		&nworkers.WorkerOptions{
			Name:       "CleanYoutubeChannelCache",
			Interval:   intervalCleanYoutubeChannelCache.DurationPtr(),
			FirstDelay: 10 * time.Second,
		},
	))

	ws.Add(nworkers.NewWorker(
		cachejobs.NewCleanCacheJob(deps.DownloadStateCache),
		&nworkers.WorkerOptions{
			Name:       "CleanDownloadStateCache",
			Interval:   intervalCleanDownloadStateCache.DurationPtr(),
			FirstDelay: 15 * time.Second,
		},
	))
}
