package wjobs

import (
	"context"
	"log/slog"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type updateHashJob struct {
	logger *slog.Logger
	runner pworkers.DownloadMaintenanceRunner
}

func NewUpdateHashJob(logger *slog.Logger, runner pworkers.DownloadMaintenanceRunner) *updateHashJob {
	return &updateHashJob{
		logger: logger,
		runner: runner,
	}
}

func (j *updateHashJob) Execute(ctx context.Context) error {
	err := j.runner.UpdateHash(ctx)
	j.logger.Debug("Job done", "name", "UpdateHash")
	return err
}
