package ytdownloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (uc *YouTubeDownloader) ScheduleDownload(
	ctx context.Context,
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
			FileName:   filename,
			YoutubeUrl: url,
		},
		options,
	)
	if err != nil {
		uc.logger.Error("Insert record failed", "error", err)
		return nil, err
	}

	file, err := uc.file.FindByFileId(ctx, fileId, true)
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

	uc.saveStateByFile(ctx, file)

	file, err = uc.addFileToQueueDownload(ctx, fileId, file.DownloadTask.TaskId)
	if err != nil {
		uc.logger.Error("Failed add to queue", "error", err)
		return nil, err
	}

	return &dto.ScheduleDownloadResponse{
		FileId:       file.FileId,
		Status:       file.Status,
		YoutubeTitle: file.YoutubeTitle,
		Format:       file.Ext,
	}, nil
}

func (uc *YouTubeDownloader) addFileToQueueDownload(ctx context.Context, fileId uuid.UUID, taskId uuid.UUID) (*ddownload.File, error) {
	err := uc.fileStatus.Pending(ctx, fileId, taskId)
	if err != nil {
		uc.logger.Debug("Failed update status", "error", err)
		uc.fileStatus.Failed(ctx, fileId, uptr.String(err.Error()))
		uc.saveStateByFileId(ctx, fileId)
		return nil, err
	}

	file, err := uc.file.FindByFileId(ctx, fileId, true)
	if err != nil {
		uc.logger.Debug("Failed find file", "error", err)
		return nil, err
	}

	uc.saveStateByFile(ctx, file)

	uc.enqueueDownloadTask(file.DownloadTask)

	return file, nil
}
