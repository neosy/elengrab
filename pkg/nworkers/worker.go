package nworkers

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

// WorkerJob represents a single unit of work to be executed by a Worker
type WorkerJob interface {
	Execute(ctx context.Context) error
}

// Worker defines the interface for a worker that can run jobs periodically.
type Worker interface {
	// Starts the worker loop
	Run(ctx context.Context, stop <-chan struct{}) error
	// Returns worker name
	Name() string
}

// worker is a concrete implementation of Worker.
type worker struct {
	// indicates if the worker is currently running
	running atomic.Bool
	// the job to execute
	job WorkerJob

	// options
	optons WorkerOptions
}

// NewWorker creates a new Worker with the given job and options.
func NewWorker(job WorkerJob, options *WorkerOptions) Worker {
	w := &worker{
		job: job,
	}

	if options != nil {
		w.optons = *options
	}

	return w
}

// Run starts the worker loop. It executes the job after FirstDelay,
// then either once or repeatedly based on Interval until stop or ctx cancellation.
func (w *worker) Run(ctx context.Context, stop <-chan struct{}) error {
	if !w.running.CompareAndSwap(false, true) {
		return errors.New("worker already running")
	}
	defer w.running.Store(false)

	// Timer for first execution after delay
	timer := time.NewTimer(w.optons.FirstDelay)
	defer timer.Stop()

	if w.optons.Interval == nil {
		// Single execution after FirstDelay
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
		// Periodic execution with Interval
		ticker := time.NewTicker(*w.optons.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-stop:
				return nil
			case <-timer.C: // first execution
				w.job.Execute(ctx)
			case <-ticker.C: // subsequent executions
				w.job.Execute(ctx)
			}
		}
	}
}

// Name returns the worker's name.
func (w *worker) Name() string {
	return w.optons.Name
}
