package workers

import (
	"time"

	ytdownloader "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	"github.com/neosy/elengrab/pkg/nworkers"
)

const (
	intervalUpdateHashDefault       = 8 * time.Hour
	intervalDeleteDuplicatesDefault = 1 * time.Hour
)

type Dependencies struct {
	// usecases
	Downloader *ytdownloader.YouTubeDownloader

	// options
	IntervalUpdateHash       time.Duration
	IntervalDeleteDuplicates time.Duration
}

func InitWorkers(ws *nworkers.Workers, deps *Dependencies) {
	var (
		intervalUpdateHash       = nworkers.NewInterval(intervalUpdateHashDefault, deps.IntervalUpdateHash)
		intervalDeleteDuplicates = nworkers.NewInterval(intervalDeleteDuplicatesDefault, deps.IntervalDeleteDuplicates)
	)

	ws.Add(nworkers.NewWorker(
		wjobs.NewUpdateHashJob(deps.Downloader),
		&nworkers.WorkerOptions{
			Name:       "UpdateHash",
			Interval:   intervalUpdateHash.DurationPtr(),
			FirstDelay: 5 * time.Second,
		},
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteDuplicatesJob(deps.Downloader),
		&nworkers.WorkerOptions{
			Name:       "DeleteDuplicates",
			Interval:   intervalDeleteDuplicates.DurationPtr(),
			FirstDelay: 30 * time.Second,
		},
	))
}
