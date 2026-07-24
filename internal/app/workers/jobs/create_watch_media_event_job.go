package wjobs

import (
	"context"
	"fmt"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type createWatchMediaEventJob struct {
	request *dto.CreateMediaWatchEventRequest
	runner  pworkers.WatchEventRunner
}

func NewCreateWatchMediaEventJob(runner pworkers.WatchEventRunner, request *dto.CreateMediaWatchEventRequest) nworkerpool.Job {
	return &createWatchMediaEventJob{
		request: request,
		runner:  runner,
	}
}

func (j *createWatchMediaEventJob) Execute(ctx context.Context, workerId uint64) error {
	return j.runner.ExecuteCreateMediaWatchEvent(ctx, workerId, j.request)
}

func (j *createWatchMediaEventJob) ID() string {
	if j.request.Event == nil {
		return ""
	}
	return j.request.Event.EventID.String()
}

func (j *createWatchMediaEventJob) Name() string {
	var id string
	if j.request.Event != nil {
		id = j.request.Event.DownloadID.String()
	}
	return fmt.Sprintf("Create media watch event for download ID: %s", id)
}
