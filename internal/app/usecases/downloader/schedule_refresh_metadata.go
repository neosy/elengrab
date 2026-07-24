package downloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/authz"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
)

func (uc *Downloader) ScheduleRefreshMetadata(
	ctx context.Context,
	authCtx dauth.UserContext,
	downloadID uuid.UUID,
) error {
	if uc.demoMode {
		uc.broadcastNotification(
			authCtx.EventKey(),
			dto.BroadcastNotificationModuleGrabForm,
			dto.BroadcastNotificationTypeError,
			"Operation not allowed in demo mode",
		)
		return exceptions.DEMO_MODE_RESTRICTION.NewErrorx()
	}

	if uc.appMode == dtypes.AppModePublicReadonly && authz.IsAnonymous(authCtx.RoleIDs) {
		uc.broadcastNotification(
			authCtx.EventKey(),
			dto.BroadcastNotificationModuleGrabForm,
			dto.BroadcastNotificationTypeError,
			"You must be authenticated to perform this action",
		)
		return exceptions.UNAUTHORIZED.NewErrorx()
	}

	err := uc.validateWriteOperation(authCtx)
	if err != nil {
		return err
	}

	download, err := uc.download.GetByDownloadIDNoCache(ctx, downloadID)
	if err != nil {
		return err
	}

	err = uc.validateDownloadWriteAccess(authCtx, download)
	if err != nil {
		return err
	}

	err = uc.downloadStatus.Refreshing(ctx, downloadID)
	if err != nil {
		return err
	}

	uc.broadcastDownloadStartRefreshing(ctx, downloadID)

	err = uc.addTaskToQueueRefreshMetadata(authCtx.UserID, download)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Downloader) addTaskToQueueRefreshMetadata(userID uuid.UUID, download *ddownload.MediaDownload) error {
	task := &ddownload.RefreshMetadataTask{
		TaskID:     uuid.New(),
		DownloadID: download.DownloadID,
		UserID:     userID,
	}

	job := uc.enqueueRefreshMetadataTask(task)
	if job == nil {
		uc.logger.Warn(
			"Task has not been added to the queue",
			"name", "Refresh metadata",
			"downloadID", download.DownloadID,
		)
		return errorx.Errorf(
			"task has not been added to the queue",
			exceptions.QUEUE_PUBLISH_FAILED,
			errorx.WithErrorMessage("Unable to schedule metadata refresh task"),
		)
	}

	return nil
}

func (uc *Downloader) enqueueRefreshMetadataTask(task *ddownload.RefreshMetadataTask) nworkerpool.Job {
	job := wjobs.NewRefreshMetadataJob(uc, task)

	if !uc.operationDispatcher.AddJob(job) {
		return nil
	}

	return job
}
