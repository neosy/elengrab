package mediawatch

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
)

func (uc *MediaWatch) CreateMediaWatchEvent(req *dto.TrackMediaWatchEventRequest, mediaDuration time.Duration) error {
	if req == nil {
		return apperrors.ErrFuncParamNullPointer
	}

	if mediaDuration == 0 {
		return errorx.NewHTTPMessage("Media duration is zero", http.StatusConflict)
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
		MediaDuration: mediaDuration,
	}

	job := uc.enqueueCreateMediaWatchEvent(newReq)
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

func (uc *MediaWatch) enqueueCreateMediaWatchEvent(req *dto.CreateMediaWatchEventRequest) nworkerpool.Job {
	job := wjobs.NewCreateWatchMediaEventJob(uc, req)

	if !uc.watchEventDispatcher.AddJob(job) {
		return nil
	}

	return job
}
