package wjobs

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewDownloadJob(runner pworkers.DownloadTaskRunner, task *ddownload.DownloadTask) workerpool.Job {
	var (
		jobID   string
		jobName string
	)

	if task.JobID != nil {
		jobID = task.JobID.String()
	} else {
		jobID = uuid.New().String()
	}

	jobName = fmt.Sprintf("Media Download: %s", task.MediaUrl)

	handler := func(ctx context.Context, workerID uint64) error {
		return runner.ExecuteDownloadTask(ctx, workerID, task)
	}

	return workerpool.NewSimpleJob(jobID, jobName, handler)
}
