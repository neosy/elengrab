package wjobs

import (
	"context"
	"log/slog"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type flushWALJob struct {
	logger *slog.Logger
	runner pworkers.MaintenanceRunner
}

func NewFlushWALJob(logger *slog.Logger, runner pworkers.MaintenanceRunner) *flushWALJob {
	return &flushWALJob{
		logger: logger,
		runner: runner,
	}
}

func (j *flushWALJob) Execute(ctx context.Context) error {
	err := j.runner.FlushWAL()
	j.logger.Debug("Job done", "name", "FlushWAL")
	return err
}
