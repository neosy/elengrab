package downloader

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	wjobs "github.com/neosy/elengrab/internal/app/workers/jobs"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/pkg/nworkerpool"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (uc *YouTubeDownloader) ScheduleDownload(
	ctx context.Context,
	userID uuid.UUID,
	url string,
	options *ddownload.DownloadOptions,
) (*dto.ScheduleDownloadResponse, error) {
	fileId := uuid.New()
	filename := fileId.String()

	options.Filename = &filename

	err := uc.file.Create(
		ctx,
		&ddownload.File{
			FileId:     fileId,
			UserID:     &userID,
			FileName:   filename,
			YoutubeUrl: url,
		},
		options,
	)
	if err != nil {
		uc.logger.Error("Insert record failed", "error", err)
		return nil, err
	}

	var accessByUserID *uuid.UUID
	if uc.historyMode != dtypes.HistoryModeGlobal {
		accessByUserID = &userID
	}

	file, err := uc.file.GetByFileId(ctx, accessByUserID, fileId)
	if err != nil {
		uc.logger.Error("Failed find file", "error", err)
		return nil, err
	}
	if file.DownloadTask == nil {
		file.DownloadTask, err = uc.dlTask.FindByFileId(ctx, fileId, true)
		if err != nil {
			uc.logger.Error("Failed find task", "error", err)
			return nil, err
		}
	}

	err = uc.addFileToQueueDownload(ctx, fileId, file.DownloadTask.TaskId)
	if err != nil {
		uc.logger.Error("Failed add to queue", "error", err)
		return nil, err
	}

	f, _ := uc.file.GetByFileId(ctx, accessByUserID, fileId)
	if f != nil {
		file = f
	}

	return &dto.ScheduleDownloadResponse{
		FileId:       file.FileId,
		Status:       file.Status,
		YoutubeTitle: file.YoutubeTitle,
		Format:       file.Ext,
	}, nil
}

func (uc *YouTubeDownloader) addFileToQueueDownload(ctx context.Context, fileId uuid.UUID, taskId uuid.UUID) error {
	var (
		file *ddownload.File
	)

	err := uc.file.Tx(ctx, func(ctx context.Context) error {
		err := uc.fileStatus.Pending(ctx, fileId, taskId, uuid.New())
		if err != nil {
			uc.logger.Warn("Failed update status", "fileId", fileId, "error", err)
			uc.dlStateCache.Delete(ctx, fileId)
			return err
		}

		file, err = uc.file.GetByFileId(ctx, nil, fileId)
		if err != nil {
			uc.dlStateCache.Delete(ctx, fileId)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	job := uc.enqueueDownloadTask(file.DownloadTask)
	if job == nil {
		err := fmt.Errorf("task has not been added to the queue")
		uc.logger.Warn("Task has not been added to the queue", "fileId", file.FileId)

		e := uc.fileStatus.Failed(ctx, fileId, nil, uptr.String("failed to enqueue download task"))
		if e != nil {
			uc.logger.Warn("Failed update status", "fileId", file.FileId, "error", e)
			uc.dlStateCache.Delete(ctx, fileId)
			return fmt.Errorf("%v: %w", err, e)
		}

		return err
	}

	return nil
}

func (uc *YouTubeDownloader) enqueueDownloadTask(task *ddownload.DownloadTask) nworkerpool.Job {
	job := wjobs.NewDownloadJob(task, uc)

	if !uc.dlDispetcher.AddJob(job) {
		return nil
	}

	return job
}
