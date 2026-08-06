package wjobs

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type refreshMetadataJob struct {
	jobID  uuid.UUID
	task   *ddownload.RefreshMetadataTask
	runner pworkers.DownloadTaskRunner
}

func NewRefreshMetadataJob(runner pworkers.DownloadTaskRunner, task *ddownload.RefreshMetadataTask) workerpool.Job {
	var jobID uuid.UUID
	if task.JobID != nil {
		jobID = *task.JobID
	}

	return &refreshMetadataJob{
		jobID:  jobID,
		task:   task,
		runner: runner,
	}
}

func (j *refreshMetadataJob) Execute(ctx context.Context, workerId uint64) error {
	return j.runner.ExecuteRefreshMetadataTask(ctx, workerId, j.task)
}

func (j *refreshMetadataJob) ID() string {
	return j.jobID.String()
}

func (j *refreshMetadataJob) Name() string {
	return fmt.Sprintf("Refresh media metadata for download ID: %s", j.task.DownloadID)
}
