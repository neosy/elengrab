package downloader

import (
	"context"
	"strings"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/app/usecases/types"
	wjobs "github.com/neosy/elengrab/internal/app/workers/pool_jobs"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func (uc *downloader) ScheduleUpload(
	ctx context.Context,
	authCtx dauth.AuthContext,
	upload pservices.UploadedFile,
	req dto.ScheduleUploadRequest,
) (*dto.ScheduleUploadResponse, error) {
	if upload == nil {
		return nil, apperrors.ErrFuncParamNullPointer
	}

	mediaInfo := upload.MediaInfo()

	if mediaInfo == nil {
		return nil, errorx.New("uploaded file has no media info",
			errorx.WithErrorMessage("Could not read information about the uploaded file."))
	}

	if err := uc.validateWriteAccess(authCtx, types.BroadcastNotificationModuleUploadFile); err != nil {
		return nil, err
	}

	broadcastError := uc.newErrorBroadcaster(authCtx, types.BroadcastNotificationModuleUploadFile)

	downloadID := uuid.New()

	var mediaDescription *string
	if desc := strings.TrimSpace(req.Description); desc != "" {
		mediaDescription = &desc
	}

	err := uc.download.CreateUpload(
		ctx,
		&ddownload.MediaDownload{
			DownloadID:               downloadID,
			SourceType:               dtypes.MediaDownloadSourceTypeUpload,
			UserID:                   &authCtx.UserID,
			Visibility:               req.Visibility,
			MediaTitleOriginal:       req.Title,
			MediaTitle:               req.Title,
			MediaDescriptionOriginal: mediaDescription,
			MediaDescription:         mediaDescription,
			Ext:                      mediaInfo.Format.Ext(),
			MediaInfo:                mediaInfo,
		},
		upload.FilePath(),
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
	if download.ProcessingTask == nil {
		download.ProcessingTask, err = uc.processingTask.GetByDownloadID(ctx, downloadID)
		if err != nil {
			uc.logger.Error("Failed find task", "error", err)
			return nil, broadcastError(err)
		}
	}

	err = uc.addUploadToQueueProcessing(ctx, downloadID, download.ProcessingTask.TaskID)
	if err != nil {
		uc.logger.Error("Failed add to queue", "error", err)
		return nil, broadcastError(err)
	}

	downloadInfo, _ := uc.resolveActualDownloadInfoByDownload(ctx, download)

	uc.broadcastMediaDownloadAdded(downloadInfo)

	return &dto.ScheduleUploadResponse{
		DownloadID: download.DownloadID,
		Status:     download.Status,
		MediaTitle: download.MediaTitle,
		MediaInfo:  download.MediaInfo,
	}, nil
}

func (uc *downloader) addUploadToQueueProcessing(ctx context.Context, downloadID uuid.UUID, taskId uuid.UUID) error {
	var (
		download *ddownload.MediaDownload
	)

	err := uc.download.Tx(ctx, func(ctx context.Context) error {
		jobID := uuid.New()
		err := uc.uploadStatus.Pending(ctx, downloadID, taskId, jobID)
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

	job := uc.enqueueUploadTask(download.ProcessingTask)
	if job == nil {
		uc.logger.Warn("Task has not been added to the queue", "downloadID", download.DownloadID)

		e := uc.uploadStatus.Failed(ctx, downloadID, nil, "failed to enqueue download task")
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

func (uc *downloader) enqueueUploadTask(task *ddownload.MediaProcessingTask) workerpool.Job {
	job := wjobs.NewUploadJob(uc, task)

	if err := uc.operationDispatcher.AddJob(job); err != nil {
		uc.logger.Warn(
			"Failed to add upload job to dispatcher",
			"jobID", job.ID(),
			"error", err,
		)
		return nil
	}

	return job
}
