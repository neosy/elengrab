package workers

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	ytdownloader "github.com/neosy/elengrab/internal/app/usecases/youtube_downloader"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
	"github.com/neosy/elengrab/pkg/worker"
)

type Dependencies struct {
	// usecases
	Downloader *ytdownloader.YouTubeDownloader
}

type Workers struct {
	logger *slog.Logger

	items []worker.Worker

	running atomic.Bool
	stop    chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
}

func NewWorkers(logger *slog.Logger, deps *Dependencies) *Workers {
	ws := &Workers{
		logger: logger,
		stop:   make(chan struct{}),
	}

	ws.Add(worker.NewWorker(
		&worker.WorkerOptions{
			Name:       "UpdateHash",
			Interval:   uptr.Any(8 * time.Hour),
			FirstDelay: 5 * time.Second,
		},
		wjobs.NewUpdateHashJob(deps.Downloader),
	))

	ws.Add(worker.NewWorker(
		&worker.WorkerOptions{
			Name:       "DeleteDuplicates",
			Interval:   uptr.Any(1 * time.Hour),
			FirstDelay: 30 * time.Second,
		},
		wjobs.NewDeleteDuplicatesJob(deps.Downloader),
	))

	return ws
}

func (w *Workers) Add(worker worker.Worker) {
	w.items = append(w.items, worker)
}

func (w *Workers) StartWorkers(ctx context.Context) error {
	if !w.running.CompareAndSwap(false, true) {
		w.logger.Debug("Workers already running")
		return errors.New("workers already running")
	}

	for _, item := range w.items {
		w.wg.Go(func() {
			item.Run(ctx, w.stop)
		})
		w.logger.Debug("Worker started", "name", item.Name())
	}

	w.logger.Debug("Workers running...", "qty", len(w.items))

	return nil
}

func (w *Workers) StopWorkers() {
	if !w.running.Load() {
		return
	}

	select {
	case <-w.stop:
		return
	default:
		w.mu.Lock()
		close(w.stop)
		w.mu.Unlock()
	}

	w.Wait()

	w.logger.Debug("Workers stopped")
}

func (w *Workers) Wait() {
	w.wg.Wait()
}
