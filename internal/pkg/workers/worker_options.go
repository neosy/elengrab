package nworkers

import "time"

// WorkerOption defines a functional option for configuring a Worker.
// It is a function that modifies a WorkerOptions struct.
type WorkerOption func(*WorkerOptions)

type WorkerOptions struct {
	// Worker name
	Name string

	// StartAt defines the absolute time when the worker should start.
	// If zero == time.Time{}, the worker starts in interval.
	StartAt time.Time

	// Interval between repeated runs
	// If zero, the worker not runs repeatedly.
	Interval time.Duration

	// OneShotDelay sets the delay before the first one-shot execution.
	// nil = skip, 0 = run immediately, >0 = run after delay.
	OneShotDelay *time.Duration

	// MaxRuns limits the total number of job executions (one-shot + repeated).
	// 0 = unlimited.
	MaxRuns int
}

// DefaultWorkerOptions returns a WorkerOptions struct with default values.
// All fields are initialized to their zero values.
func DefaultWorkerOptions() WorkerOptions {
	return WorkerOptions{
		Name:         "",
		StartAt:      time.Time{},
		Interval:     0,
		OneShotDelay: nil,
		MaxRuns:      0,
	}
}

// WithName returns a WorkerOption that sets the Name field
// of WorkerOptions to the provided string.
func WithName(name string) WorkerOption {
	return func(o *WorkerOptions) {
		o.Name = name
	}
}

// WithStartAt returns a WorkerOption that sets the StartAt field
// of WorkerOptions to the specified time.Time value.
func WithStartAt(startAt time.Time) WorkerOption {
	return func(o *WorkerOptions) {
		o.StartAt = startAt
	}
}

// WithInterval returns a WorkerOption that sets the Interval field
// of WorkerOptions to the given time.Duration.
func WithInterval(interval time.Duration) WorkerOption {
	return func(o *WorkerOptions) {
		o.Interval = interval
	}
}

// WithIntervalDefault sets Interval to interval,
// or intervalDefault if interval is 0.
func WithIntervalDefault(interval, intervalDefault time.Duration) WorkerOption {
	return func(o *WorkerOptions) {
		if interval == 0 {
			o.Interval = intervalDefault
		} else {
			o.Interval = interval
		}
	}
}

// WithInitialDelay sets a delay before the first execution of the worker.
// If not set, the worker uses Interval as the initial delay.
func WithInitialDelay(delay time.Duration) WorkerOption {
	return func(o *WorkerOptions) {
		o.OneShotDelay = &delay
	}
}

// WithMaxRuns returns a WorkerOption that sets the MaxRuns field
// of WorkerOptions to the specified integer value.
func WithMaxRuns(maxRuns int) WorkerOption {
	return func(o *WorkerOptions) {
		o.MaxRuns = maxRuns
	}
}
