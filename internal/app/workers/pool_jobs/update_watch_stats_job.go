package wjobs

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewUpdateWatchStatsJob(runner pworkers.WatchEventRunner) workerpool.Job {
	var (
		jobID   string
		jobName string
	)

	jobID = uuid.New().String()
	jobName = "Update watch stats"

	handler := func(ctx context.Context, workerID uint64) error {
		return runner.ExecuteUpdateStats(ctx, workerID)
	}

	return workerpool.NewSimpleJob(jobID, jobName, handler)
}
