package wjobs

import (
	"context"
	"log/slog"
	"time"

	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type startupAuthWebJob struct {
	logger *slog.Logger
	runner pworkers.AuthWebMaintenanceRunner
}

func NewStartupAuthWebJob(logger *slog.Logger, runner pworkers.AuthWebMaintenanceRunner) *startupAuthWebJob {
	return &startupAuthWebJob{
		logger: logger,
		runner: runner,
	}
}

func (j *startupAuthWebJob) Execute(ctx context.Context) error {
	startTime := time.Now()
	err := j.runner.Startup(ctx)
	elapsed := time.Since(startTime)

	j.logger.Debug(
		"Job done",
		"name", "StartupMaintenanceAuthWeb",
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return err
}
