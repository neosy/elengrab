package nworkers

import (
	"context"
	"log/slog"
	"time"

	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

type baseJob struct {
	name    string
	execute func(ctx context.Context, job Job) error
}

type JobExecute func(ctx context.Context, job Job) error

// Job represents a single unit of work to be executed by a Worker
type Job interface {
	Execute(ctx context.Context) error
	Name() string
}

// NewJob creates a new Job with the given name and execution function.
func NewJob(name string, execute JobExecute) Job {
	return &baseJob{
		name:    name,
		execute: execute,
	}
}

// Execute runs the job's execution function with the provided context.
func (j *baseJob) Execute(ctx context.Context) error {
	err := j.execute(ctx, j)
	return err
}

// Name returns the name of the job.
func (j *baseJob) Name() string {
	return j.name
}

// WrapJobExecute wraps a job execution function with logging and optional timing measurement.
func WrapJobExecute(
	logger *slog.Logger,
	run func(ctx context.Context) error,
	opts ...JobExecuteOption,
) JobExecute {
	options := NewJobExecuteOptions(opts...)

	return func(ctx context.Context, j Job) error {
		var start time.Time
		if options.measureElapsed {
			start = time.Now()
		}

		// Execute the job
		err := run(ctx)

		// Log the job completion
		if logger != nil {
			logger = logger.With("jobName", j.Name())

			// Add elapsed time to the logger if requested
			if options.measureElapsed {
				logger = logger.With("elapsed", uformat.DurationFormat(time.Since(start)))
			}

			logger.Debug("Job done")
		}

		return err
	}
}
