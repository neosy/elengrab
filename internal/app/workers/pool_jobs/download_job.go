package wjobs

import (
	"context"
	"fmt"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewDownloadJob(runner pworkers.DownloadTaskRunner, task *ddownload.DownloadTask) workerpool.Job {
	jobID := task.JobID.String()
	jobName := fmt.Sprintf("Media Download: %s", task.MediaUrl)

	handler := func(ctx context.Context, workerID uint64) error {
		return runner.ExecuteDownloadTask(ctx, workerID, task)
	}

	return workerpool.NewSimpleJob(jobID, jobName, handler)
}
