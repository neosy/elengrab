package wjobs

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewCreateWatchMediaEventJob(
	runner pworkers.WatchEventRunner,
	request *dto.CreateMediaWatchEventRequest,
) workerpool.Job {
	var (
		jobID   string
		jobName string
	)

	if request.Event != nil {
		jobID = request.Event.EventID.String()
		jobName = request.Event.DownloadID.String()
	}

	handler := func(ctx context.Context, workerID uint64) error {
		return runner.ExecuteCreateMediaWatchEvent(ctx, workerID, request)
	}

	return workerpool.NewSimpleJob(jobID, jobName, handler)
}
