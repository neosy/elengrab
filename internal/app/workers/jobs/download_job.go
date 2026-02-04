package wjobs

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type downloadJob struct {
	uuid   uuid.UUID
	task   *ddownload.DownloadTask
	runner pworkers.DownloadTaskRunner
}

func NewDownloadJob(task *ddownload.DownloadTask, runner pworkers.DownloadTaskRunner) *downloadJob {
	var jobID uuid.UUID
	if task.JobID != nil {
		jobID = *task.JobID
	}

	return &downloadJob{
		uuid:   jobID,
		task:   task,
		runner: runner,
	}
}

func (j *downloadJob) Execute(ctx context.Context, workerId uint) error {
	return j.runner.ExecuteDownloadTask(ctx, workerId, j.task)
}

func (j *downloadJob) ID() string {
	return j.uuid.String()
}

func (j *downloadJob) Name() string {

	return fmt.Sprintf("Media Download: %s", j.task.MediaUrl)
}
