package wjobs

import (
	"context"
	"log/slog"
	"time"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
	uformat "github.com/neosy/elengrab/pkg/utils/format"
)

type startupDatabaseJob struct {
	logger *slog.Logger
	runner pworkers.MaintenanceRunner
}

func NewStartupDatabaseJob(logger *slog.Logger, runner pworkers.MaintenanceRunner) *startupDatabaseJob {
	return &startupDatabaseJob{
		logger: logger,
		runner: runner,
	}
}

func (j *startupDatabaseJob) Execute(ctx context.Context) error {
	startTime := time.Now()
	err := j.runner.StartupDatabase(ctx)
	elapsed := time.Since(startTime)

	j.logger.Debug(
		"Job done",
		"name", "StartupMaintenanceDatabase",
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return err
}
