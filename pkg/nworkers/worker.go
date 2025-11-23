package nworkers

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

type WorkerJob interface {
	Execute(ctx context.Context) error
}

type Worker interface {
	Run(ctx context.Context, stop <-chan struct{}) error
	Name() string
}

type worker struct {
	running atomic.Bool
	job     WorkerJob

	// options
	optons WorkerOptions
}

func NewWorker(job WorkerJob, options *WorkerOptions) Worker {
	w := &worker{
		job: job,
	}

	if options != nil {
		w.optons = *options
	}

	return w
}

func (w *worker) Run(ctx context.Context, stop <-chan struct{}) error {
	if !w.running.CompareAndSwap(false, true) {
		return errors.New("worker already running")
	}
	defer w.running.Store(false)

	// first start after
	timer := time.NewTimer(w.optons.FirstDelay)
	defer timer.Stop()

	if w.optons.Interval == nil {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-stop:
				return nil
			case <-timer.C:
				w.job.Execute(ctx)
				return nil
			}
		}
	} else {
		ticker := time.NewTicker(*w.optons.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-stop:
				return nil
			case <-timer.C:
				w.job.Execute(ctx)
			case <-ticker.C:
				w.job.Execute(ctx)
			}
		}
	}
}

func (w *worker) Name() string {
	return w.optons.Name
}
