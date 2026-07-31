package nworkers

import "time"

// WorkerOption defines a functional option for configuring a Worker.
// It is a function that modifies a WorkerOptions struct.
type WorkerOption func(*WorkerOptions)

type WorkerOptions struct {
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
		StartAt:      time.Time{},
		Interval:     0,
		OneShotDelay: nil,
		MaxRuns:      0,
	}
}

// ApplyWorkerOptions applies the provided WorkerOption functions to the given WorkerOptions.
func ApplyWorkerOptions(options *WorkerOptions, opts ...WorkerOption) {
	for _, opt := range opts {
		opt(options)
	}
}

// NewJobWorkerOptions creates a new WorkerOptions instance with the provided options applied.
func NewWorkerOptions(opts ...WorkerOption) WorkerOptions {
	options := DefaultWorkerOptions()

	ApplyWorkerOptions(&options, opts...)

	return options
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

// WithIntervalFallback returns a WorkerOption that sets the Interval field
// of WorkerOptions to the specified interval if it is greater than zero.
// If the provided interval is zero, it sets the Interval to the fallback value.
func WithIntervalFallback(interval, fallback time.Duration) WorkerOption {
	return func(o *WorkerOptions) {
		if interval > 0 {
			o.Interval = interval
		} else {
			o.Interval = fallback
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
