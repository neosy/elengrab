package workers

import (
	"time"

	ytdownloader "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	"github.com/neosy/elengrab/pkg/nworkers"
)

const (
	intervalUpdateHashDefault         = 8 * time.Hour
	intervalDeleteDuplicatesDefault   = 1 * time.Hour
	intervalDeleteMissingFilesDefault = 30 * time.Minute
	intervalDeleteFailedDownloads     = 1 * time.Hour
)

type Dependencies struct {
	// usecases
	Downloader *ytdownloader.YouTubeDownloader

	// options
	IntervalUpdateHash            time.Duration
	IntervalDeleteDuplicates      time.Duration
	IntervalDeleteMissingFiles    time.Duration
	IntervalDeleteFailedDownloads time.Duration
	EnableMoveUnmatchedFiles      bool
}

func InitWorkers(ws *nworkers.Workers, deps *Dependencies) {
	var (
		intervalUpdateHash            = nworkers.NewInterval(intervalUpdateHashDefault, deps.IntervalUpdateHash)
		intervalDeleteDuplicates      = nworkers.NewInterval(intervalDeleteDuplicatesDefault, deps.IntervalDeleteDuplicates)
		intervalDeleteMissingFiles    = nworkers.NewInterval(intervalDeleteMissingFilesDefault, deps.IntervalDeleteMissingFiles)
		intervalDeleteFailedDownloads = nworkers.NewInterval(intervalDeleteFailedDownloads, deps.IntervalDeleteFailedDownloads)
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
}
