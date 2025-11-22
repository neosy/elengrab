package wjobs

import (
	"context"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type updateHashJob struct {
	runner pworkers.MaintenanceRunner
}

func NewUpdateHashJob(runner pworkers.MaintenanceRunner) *updateHashJob {
	return &updateHashJob{
		runner: runner,
	}
}

func (j *updateHashJob) Execute(ctx context.Context) error {
	return j.runner.UpdateHash(ctx)
}
