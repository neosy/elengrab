package wjobs

import (
	"context"
	"log/slog"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
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
	err := j.runner.BackupDatabase(ctx)
	j.logger.Debug("Job done", "name", "BackupDatabase")
	return err
}
