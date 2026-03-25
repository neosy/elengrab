package wjobs

import (
	"context"
	"log/slog"
	"time"

	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
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
