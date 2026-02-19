package wjobs

import (
	"context"
	"log/slog"
	"time"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
	uformat "github.com/neosy/elengrab/pkg/utils/format"
)

type backupDatabaseJob struct {
	logger *slog.Logger
	runner pworkers.MaintenanceRunner
}

func NewbackupDatabaseJob(logger *slog.Logger, runner pworkers.MaintenanceRunner) *backupDatabaseJob {
	return &backupDatabaseJob{
		logger: logger,
		runner: runner,
	}
}

func (j *backupDatabaseJob) Execute(ctx context.Context) error {
	startTime := time.Now()
	err := j.runner.BackupDatabase(ctx)
	elapsed := time.Since(startTime)

	j.logger.Debug(
		"Job done",
		"name", "BackupDatabase",
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return err
}
