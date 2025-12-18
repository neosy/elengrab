package wjobs

import (
	"context"
	"log/slog"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type deleteFailedDownloadsJob struct {
	logger *slog.Logger
	runner pworkers.MaintenanceRunner
}

func NewDeleteFailedDownloadsJob(logger *slog.Logger, runner pworkers.MaintenanceRunner) *deleteFailedDownloadsJob {
	return &deleteFailedDownloadsJob{
		logger: logger,
		runner: runner,
	}
}

func (j *deleteFailedDownloadsJob) Execute(ctx context.Context) error {
	err := j.runner.DeleteFailedDownloads(ctx)
	j.logger.Debug("Job done", "name", "DeleteFailedDownloads")
	return err
}
