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
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
)

func (uc *Downloader) ScheduleDownload(
	ctx context.Context,
	authCtx dauth.UserContext,
	url string,
	options *ddownload.DownloadOptions,
) (*dto.ScheduleDownloadResponse, error) {
	if uc.demoMode {
		uc.broadcastNotification(
			authCtx.EventKey(),
			dto.BroadcastNotificationModuleGrabForm,
			dto.BroadcastNotificationTypeError,
			"Operation not allowed in demo mode",
		)
		return nil, exceptions.DEMO_MODE_RESTRICTION.NewErrorx()
	}

	if uc.appMode == dtypes.AppModePublicReadonly && authz.IsAnonymous(authCtx.RoleIDs) {
		uc.broadcastNotification(
			authCtx.EventKey(),
			dto.BroadcastNotificationModuleGrabForm,
			dto.BroadcastNotificationTypeError,
			"You must be authenticated to perform this action",
		)
		return nil, exceptions.UNAUTHORIZED.NewErrorx()
	}

	returnErr := func(err error) error {
		uc.broadcastNotification(
			authCtx.EventKey(),
			dto.BroadcastNotificationModuleGrabForm,
			dto.BroadcastNotificationTypeError,
			err.Error(),
		)
		return err
	}

	downloadID := uuid.New()
	filename := downloadID.String()

	options.Filename = &filename

	var mediaVisibility dtypes.MediaVisibility
	if authCtx.UserID != uuid.Nil {
		mediaVisibility = dtypes.MediaVisibilityAuthenticated
	}

	err := uc.download.Create(
		ctx,
		&ddownload.MediaDownload{
			DownloadID: downloadID,
			UserID:     &authCtx.UserID,
			FileName:   filename,
			MediaURL:   url,
			Visibility: mediaVisibility,
		},
		options,
	)
	if err != nil {
		uc.logger.Error("Insert record failed", "error", err)
		return nil, returnErr(err)
	}

	var accessByUserID *uuid.UUID
	if uc.authz.ShouldRestrictDownloads(authCtx.RoleIDs) {
		accessByUserID = &authCtx.UserID
	}

	download, err := uc.download.GetByDownloadID(ctx, accessByUserID, downloadID)
	if err != nil {
		uc.logger.Error("Failed find download", "error", err)
		return nil, returnErr(err)
	}
	if download.DownloadTask == nil {
		download.DownloadTask, err = uc.dlTask.GetByDownloadID(ctx, downloadID)
		if err != nil {
			uc.logger.Error("Failed find task", "error", err)
			return nil, returnErr(err)
		}
	}

	err = uc.addDownloadToQueueDownload(ctx, downloadID, download.DownloadTask.TaskID)
	if err != nil {
		uc.logger.Error("Failed add to queue", "error", err)
		return nil, returnErr(err)
	}

	tmpDownload, _ := uc.download.GetByDownloadID(ctx, accessByUserID, downloadID)
	if tmpDownload != nil {
		download = tmpDownload
	}

	uc.broadcastDownloadAdd(download)

	return &dto.ScheduleDownloadResponse{
		URL:        url,
		DownloadID: download.DownloadID,
		Status:     download.Status,
		MediaTitle: download.MediaTitle,
		Format:     download.Ext,
	}, nil
}

func (uc *Downloader) addDownloadToQueueDownload(ctx context.Context, downloadID uuid.UUID, taskId uuid.UUID) error {
	var (
		download *ddownload.MediaDownload
	)

	err := uc.download.Tx(ctx, func(ctx context.Context) error {
		jobID := uuid.New()
		err := uc.downloadStatus.Pending(ctx, downloadID, taskId, jobID)
		if err != nil {
			uc.logger.Warn("Failed update status", "downloadID", downloadID, "error", err)
			uc.dlStateCache.Delete(ctx, downloadID)
			return err
		}

		download, err = uc.download.GetByDownloadID(ctx, nil, downloadID)
		if err != nil {
			uc.dlStateCache.Delete(ctx, downloadID)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	job := uc.enqueueDownloadTask(download.DownloadTask)
	if job == nil {
		uc.logger.Warn("Task has not been added to the queue", "downloadID", download.DownloadID)

		e := uc.downloadStatus.Failed(ctx, downloadID, nil, uptr.String("failed to enqueue download task"))
		if e != nil {
			uc.logger.Warn("Failed update status", "downloadID", download.DownloadID, "error", e)
			uc.dlStateCache.Delete(ctx, downloadID)
			return errorx.Errorf(
				"task has not been added to the queue: %w", e,
				exceptions.QUEUE_PUBLISH_FAILED,
				errorx.WithErrorMessage("We couldn’t process your request"),
			)
		}

		return err
	}

	return nil
}

func (uc *Downloader) enqueueDownloadTask(task *ddownload.DownloadTask) nworkerpool.Job {
	job := wjobs.NewDownloadJob(uc, task)

	if !uc.dlDispetcher.AddJob(job) {
		return nil
	}

	return job
}
