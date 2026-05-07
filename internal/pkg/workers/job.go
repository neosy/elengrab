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

func NewJob(name string, execute JobExecute) Job {
	return &baseJob{
		name:    name,
		execute: execute,
	}
}

func (j *baseJob) Execute(ctx context.Context) error {
	return j.execute(ctx, j)
}

func (j *baseJob) Name() string {
	return j.name
}

// MakeTimedJobExecute wraps job execution with elapsed time measurement
// and debug logging.
func MakeTimedJobExecute(logger *slog.Logger, run func(ctx context.Context) error) JobExecute {
	return func(ctx context.Context, j Job) error {
		start := time.Now()
		err := run(ctx)

		logger.Debug(
			"Job done",
			"name", j.Name(),
			"elapsed", uformat.DurationFormat(time.Since(start)),
		)

		return err
	}
}
