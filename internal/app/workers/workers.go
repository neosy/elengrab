package workers

import (
	"time"

	ytdownloader "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	"github.com/neosy/elengrab/pkg/nworkers"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

type Dependencies struct {
	// usecases
	Downloader *ytdownloader.YouTubeDownloader
}

func InitWorkers(ws *nworkers.Workers, deps *Dependencies) {
	ws.Add(nworkers.NewWorker(
		wjobs.NewUpdateHashJob(deps.Downloader),
		&nworkers.WorkerOptions{
			Name:       "UpdateHash",
			Interval:   uptr.Any(8 * time.Hour),
			FirstDelay: 5 * time.Second,
		},
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteDuplicatesJob(deps.Downloader),
		&nworkers.WorkerOptions{
			Name:       "DeleteDuplicates",
			Interval:   uptr.Any(1 * time.Hour),
			FirstDelay: 30 * time.Second,
		},
	))
}
