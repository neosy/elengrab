package wjobs

import (
	"context"

	"github.com/neosy/elengrab/internal/pkg/workerpool"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type updateWatchStatsJob struct {
	runner pworkers.WatchEventRunner
}

func NewUpdateWatchStatsJob(runner pworkers.WatchEventRunner) workerpool.Job {
	return &updateWatchStatsJob{
		runner: runner,
	}
}

func (j *updateWatchStatsJob) Execute(ctx context.Context, workerId uint64) error {
	return j.runner.ExecuteUpdateStats(ctx, workerId)
}

func (j *updateWatchStatsJob) ID() string {
	return ""
}

func (j *updateWatchStatsJob) Name() string {
	return "Update watch stats"
}
