package mediawatch

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	wjobs "github.com/neosy/elengrab/internal/app/workers/pool_jobs"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func (uc *MediaWatch) CreateMediaWatchEvent(
	req *dto.TrackMediaWatchEventRequest,
	mediaDuration time.Duration,
	runner pworkers.DownloadTaskRunner,
) error {
	if req == nil {
		return apperrors.ErrFuncParamNullPointer
	}

	if mediaDuration == 0 {
		return errorx.NewHTTPMessage("Media duration is zero", http.StatusConflict)
	}

	req.AdjustForMediaDuration(mediaDuration)

	if err := req.Validate(); err != nil {
		return errorx.Errorf(
			"invalid track media watch event request: %w",
			err, errorx.WithHttpStatus(http.StatusBadRequest),
		)
	}

	event := &ddownload.MediaWatchEvent{
		EventID:    uuid.New(),
		DownloadID: req.DownloadID,
		UserID:     req.UserID,
		SessionID:  req.SessionID,
		Position:   req.Position,
		Interval:   req.Interval,
	}

	newReq := &dto.CreateMediaWatchEventRequest{
		Event:         event,
		EventType:     req.EventType,
		MediaDuration: mediaDuration,
	}

	job := uc.enqueueCreateMediaWatchEvent(newReq, runner)
	if job == nil {
		uc.logger.Warn(
			"Task has not been added to the queue",
			"name", "Watch event",
			"downloadID", event.DownloadID,
		)
		return errorx.Errorf(
			"failed to enqueue media watch event",
			exceptions.QUEUE_PUBLISH_FAILED,
			errorx.WithErrorMessage("Unable to enqueue media watch event"),
		)
	}

	return nil
}

func (uc *MediaWatch) enqueueCreateMediaWatchEvent(
	req *dto.CreateMediaWatchEventRequest,
	runner pworkers.DownloadTaskRunner,
) workerpool.Job {
	job := wjobs.NewCreateWatchMediaEventJob(runner, req)

	if err := uc.watchEventDispatcher.AddJob(job); err != nil {
		uc.logger.Error(
			"Failed to enqueue create media watch event job",
			"eventID", req.Event.EventID,
			"error", err,
		)
		return nil
	}

	return job
}
