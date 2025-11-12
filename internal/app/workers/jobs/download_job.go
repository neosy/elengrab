package wjobs

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type downloadJob struct {
	task   *ddownload.DownloadTask
	runner pworkers.DownloadTaskRunner
}

func NewDownloadJob(task *ddownload.DownloadTask, runner pworkers.DownloadTaskRunner) *downloadJob {
	return &downloadJob{
		task:   task,
		runner: runner,
	}
}

func (j *downloadJob) Execute(ctx context.Context, workerId uint) error {
	return j.runner.ExecuteDownloadTask(ctx, workerId, j.task)
}
