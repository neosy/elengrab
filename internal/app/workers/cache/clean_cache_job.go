package cachejobs

import (
	"context"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type cleanCacheJob struct {
	runner pworkers.CacheRunner
}

func NewCleanCacheJob(runner pworkers.CacheRunner) *cleanCacheJob {
	return &cleanCacheJob{
		runner: runner,
	}
}

func (j *cleanCacheJob) Execute(ctx context.Context) error {
	return j.runner.CleanExpired(ctx)
}
