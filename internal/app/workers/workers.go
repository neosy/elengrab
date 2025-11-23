package workers

import (
	"time"

	ytdownloader "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	"github.com/neosy/elengrab/pkg/nworkers"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

const (
	intervalUpdateHashDefault       = 8 * time.Hour
	intervalDeleteDuplicatesDefault = 1 * time.Hour
)

type interval struct {
	vDef time.Duration
	v    time.Duration
}

func newInterval(valueDefault time.Duration, value time.Duration) interval {
	return interval{
		vDef: valueDefault,
		v:    value,
	}
}

func (i *interval) value() time.Duration {
	v := i.vDef
	if i.v.Seconds() != 0 {
		v = i.v
	}
	return v
}

func (i *interval) valuePtr() *time.Duration {
	return uptr.Any(i.value())
}

type Dependencies struct {
	// usecases
	Downloader *ytdownloader.YouTubeDownloader

	// options
	IntervalUpdateHash       time.Duration
	IntervalDeleteDuplicates time.Duration
}

func InitWorkers(ws *nworkers.Workers, deps *Dependencies) {
	var (
		intervalUpdateHash       = newInterval(intervalUpdateHashDefault, deps.IntervalUpdateHash)
		intervalDeleteDuplicates = newInterval(intervalDeleteDuplicatesDefault, deps.IntervalDeleteDuplicates)
	)

	ws.Add(nworkers.NewWorker(
		wjobs.NewUpdateHashJob(deps.Downloader),
		&nworkers.WorkerOptions{
			Name:       "UpdateHash",
			Interval:   intervalUpdateHash.valuePtr(),
			FirstDelay: 5 * time.Second,
		},
	))

	ws.Add(nworkers.NewWorker(
		wjobs.NewDeleteDuplicatesJob(deps.Downloader),
		&nworkers.WorkerOptions{
			Name:       "DeleteDuplicates",
			Interval:   intervalDeleteDuplicates.valuePtr(),
			FirstDelay: 30 * time.Second,
		},
	))
}
