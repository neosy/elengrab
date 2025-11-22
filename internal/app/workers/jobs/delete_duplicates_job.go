package wjobs

import (
	"context"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type deleteDuplicatesJob struct {
	runner pworkers.MaintenanceRunner
}

func NewDeleteDuplicatesJob(runner pworkers.MaintenanceRunner) *deleteDuplicatesJob {
	return &deleteDuplicatesJob{
		runner: runner,
	}
}

func (j *deleteDuplicatesJob) Execute(ctx context.Context) error {
	return j.runner.DeleteDuplicates(ctx)
}
