package downloader

import (
	"context"

	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/app/usecases/types"
	wjobs "github.com/neosy/elengrab/internal/app/workers/pool_jobs"
	iconfig "github.com/neosy/elengrab/internal/config"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
)

func (uc *downloader) ScheduleDownload(
	ctx context.Context,
	authCtx dauth.AuthContext,
	url string,
	options *ddownload.DownloadOptions,
) (*dto.ScheduleDownloadResponse, error) {
	if err := uc.validateWriteAccess(authCtx, types.BroadcastNotificationModuleGrabForm); err != nil {
		return nil, err
	}

	broadcastError := uc.newErrorBroadcaster(authCtx, types.BroadcastNotificationModuleGrabForm)

	downloadID := uuid.New()
	uniqueFileNameNoExt := downloadID.String()

	options.UniqueFileName = &uniqueFileNameNoExt

	mediaVisibility := iconfig.InitialMediaVisibility()
	if authCtx.UserID == uuid.Nil() {
		mediaVisibility = dtypes.MediaVisibilityPublic
	}

	err := uc.download.CreateDownload(
		ctx,
		&ddownload.MediaDownload{
			DownloadID:    downloadID,
			SourceType:    dtypes.MediaDownloadSourceTypeURL,
			UserID:        &authCtx.UserID,
			FileNameNoExt: uniqueFileNameNoExt,
			MediaURL:      &url,
			Visibility:    mediaVisibility,
		},
		options,
	)
	if err != nil {
		uc.logger.Error("Insert record failed", "error", err)
		return nil, broadcastError(err)
	}

	download, err := uc.download.GetByDownloadIDNoCache(ctx, downloadID)
	if err != nil {
		uc.logger.Error("Failed find download", "error", err)
		return nil, broadcastError(err)
	}
	if download.DownloadTask == nil {
		download.DownloadTask, err = uc.dlTask.GetByDownloadID(ctx, downloadID)
		if err != nil {
			uc.logger.Error("Failed find task", "error", err)
			return nil, broadcastError(err)
		}
	}

	err = uc.addDownloadToQueueDownload(ctx, downloadID, download.DownloadTask.TaskID)
	if err != nil {
		uc.logger.Error("Failed add to queue", "error", err)
		return nil, broadcastError(err)
	}

	tmpDownload, _ := uc.download.GetByDownloadIDNoCache(ctx, downloadID)
	if tmpDownload != nil {
		download = tmpDownload
	}

	downloadInfo, _ := uc.resolveActualDownloadInfoByDownload(ctx, download)

	uc.broadcastMediaDownloadAdded(downloadInfo)

	return &dto.ScheduleDownloadResponse{
		URL:        url,
		DownloadID: download.DownloadID,
		Status:     download.Status,
		MediaTitle: download.MediaTitle,
		Format:     download.FileExtNoDot,
	}, nil
}

func (uc *downloader) addDownloadToQueueDownload(ctx context.Context, downloadID uuid.UUID, taskId uuid.UUID) error {
	var (
		download *ddownload.MediaDownload
	)

	err := uc.download.Tx(ctx, func(ctx context.Context) error {
		jobID := uuid.New()
		err := uc.downloadStatus.Pending(ctx, downloadID, taskId, jobID)
		if err != nil {
			uc.logger.Warn("Failed update status", "downloadID", downloadID, "error", err)
			return err
		}

		download, err = uc.download.GetByDownloadIDNoCache(ctx, downloadID)
		if err != nil {
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

		e := uc.downloadStatus.Failed(ctx, downloadID, nil, "failed to enqueue download task")
		if e != nil {
			uc.logger.Warn("Failed update status", "downloadID", download.DownloadID, "error", e)

			return errorx.Errorf(
				"task has not been added to the queue: %w", e,
				exceptions.QUEUE_PUBLISH_FAILED,
				errorx.WithErrorMessage("We couldn’t process your request"),
			)
		}

		return errorx.Errorf(
			"task has not been added to the queue",
			exceptions.QUEUE_PUBLISH_FAILED,
			errorx.WithErrorMessage("We couldn’t process your request"),
		)
	}

	return nil
}

func (uc *downloader) enqueueDownloadTask(task *ddownload.DownloadTask) workerpool.Job {
	job := wjobs.NewDownloadJob(uc, task)

	if err := uc.downloadDispatcher.AddJob(job); err != nil {
		uc.logger.Warn(
			"Failed to add download job to dispatcher",
			"jobID", job.ID(),
			"error", err,
		)
		return nil
	}

	return job
}
