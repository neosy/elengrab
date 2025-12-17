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

// Run starts the worker execution loop.
//
// Behavior:
//   - If RunImmediatelyDelay is set, the job is executed once after the given delay.
//   - If Interval is set, the job is executed repeatedly with the given interval.
//   - If both options are set, the first execution happens after RunImmediatelyDelay,
//     and subsequent executions happen on Interval.
//   - If neither option is set, Run returns immediately.
//
// The worker stops when:
//   - the context is cancelled,
//   - the stop channel is closed or receives a value.
func (w *worker) Run(ctx context.Context, stop <-chan struct{}) error {
	if !w.running.CompareAndSwap(false, true) {
		return errors.New("worker already running")
	}
	defer w.running.Store(false)

	var (
		timer   *time.Timer
		timerC  <-chan time.Time
		ticker  *time.Ticker
		tickerC <-chan time.Time
	)

	// One-shot execution after an optional initial delay.
	if w.optons.OneShotDelay != nil {
		timer = time.NewTimer(*w.optons.OneShotDelay)
		defer timer.Stop()
		timerC = timer.C
	} else {
		timerC = nil
	}

	// Periodic execution
	if w.optons.Interval != nil {
		ticker = time.NewTicker(*w.optons.Interval)
		defer ticker.Stop()
		tickerC = ticker.C
	} else {
		tickerC = nil
	}

	// Nothing to execute: neither delayed nor periodic execution is configured.
	if timerC == nil && tickerC == nil {
		return nil
	}

	runs := 0
	for {
		select {
		case <-ctx.Done():
			// Context cancellation always terminates the worker
			return nil
		case <-stop:
			// Explicit stop signal terminates the worker.
			return nil
		case <-timerC:
			// Initial one-shot execution.
			timerC = nil
			w.job.Execute(ctx)
			runs++
			if tickerC == nil {
				return nil
			}
			if w.optons.MaxRuns > 0 && runs >= w.optons.MaxRuns {
				return nil
			}
		case <-tickerC:
			// Periodic execution.
			w.job.Execute(ctx)
			runs++
			if w.optons.MaxRuns > 0 && runs >= w.optons.MaxRuns {
				return nil
			}
		}
	}
}

// Name returns the worker's name.
func (w *worker) Name() string {
	return w.optons.Name
}
