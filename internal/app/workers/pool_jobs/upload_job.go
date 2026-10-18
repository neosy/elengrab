package wjobs

import (
	"context"
	"fmt"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewUploadJob(runner pworkers.UploadTaskRunner, task *ddownload.MediaProcessingTask) workerpool.Job {
	jobID := task.JobID.String()
	jobName := fmt.Sprintf("Media upload: %s", task.SourcePath)

	handler := func(ctx context.Context, workerID uint64) error {
		return runner.ExecuteUploadTask(ctx, workerID, task)
	}

	return workerpool.NewSimpleJob(jobID, jobName, handler)
}
