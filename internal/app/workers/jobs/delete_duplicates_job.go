package wjobs

import (
	"context"
	"log/slog"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type deleteDuplicatesJob struct {
	logger *slog.Logger
	runner pworkers.MaintenanceRunner
}

func NewDeleteDuplicatesJob(logger *slog.Logger, runner pworkers.MaintenanceRunner) *deleteDuplicatesJob {
	return &deleteDuplicatesJob{
		logger: logger,
		runner: runner,
	}
}

func (j *deleteDuplicatesJob) Execute(ctx context.Context) error {
	err := j.runner.DeleteDuplicates(ctx)
	j.logger.Debug("Job done", "name", "DeleteDuplicates")
	return err
}
