package wjobs

import (
	"context"
	"log/slog"
	"time"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
	uformat "github.com/neosy/elengrab/pkg/utils/format"
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
	startTime := time.Now()
	err := j.runner.FlushWAL()
	elapsed := time.Since(startTime)

	j.logger.Debug(
		"Job done",
		"name", "FlushWAL",
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return err
}
