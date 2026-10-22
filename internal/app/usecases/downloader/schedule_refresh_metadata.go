package downloader

import (
	"context"

	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/types"
	wjobs "github.com/neosy/elengrab/internal/app/workers/pool_jobs"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
)

func (uc *downloader) ScheduleRefreshMetadata(
	ctx context.Context,
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
) error {
	if err := uc.validateWriteAccess(authCtx, types.BroadcastNotificationModuleResultRow); err != nil {
		return err
	}

	download, err := uc.download.GetByDownloadIDNoCache(ctx, downloadID)
	if err != nil {
		return err
	}

	if err := uc.validateDownloadEditAccess(authCtx, download); err != nil {
		return err
	}

	err = uc.downloadStatus.Refreshing(ctx, downloadID)
	if err != nil {
		return err
	}

	uc.broadcastMediaDownloadStartRefreshing(ctx, downloadID)

	err = uc.addTaskToQueueRefreshMetadata(authCtx.UserID, download)
	if err != nil {
		return err
	}

	return nil
}

func (uc *downloader) addTaskToQueueRefreshMetadata(userID uuid.UUID, download *ddownload.MediaDownload) error {
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

func (uc *downloader) enqueueRefreshMetadataTask(task *ddownload.RefreshMetadataTask) workerpool.Job {
	job := wjobs.NewRefreshMetadataJob(uc, task)

	if err := uc.operationDispatcher.AddJob(job); err != nil {
		uc.logger.Warn(
			"Failed to enqueue refresh metadata task",
			"taskID", task.TaskID,
			"error", err,
		)
		return nil
	}

	return job
}
