package wjobs

import (
	"context"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type deleteFailedDownloadsJob struct {
	runner pworkers.MaintenanceRunner
}

func NewDeleteFailedDownloadsJob(runner pworkers.MaintenanceRunner) *deleteFailedDownloadsJob {
	return &deleteFailedDownloadsJob{
		runner: runner,
	}
}

func (j *deleteFailedDownloadsJob) Execute(ctx context.Context) error {
	return j.runner.DeleteFailedDownloads(ctx)
}
