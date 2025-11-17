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
		return nil, err
	}

	uc.saveStateByFile(ctx, file)

	err = uc.fileStatus.Pending(ctx, fileId)
	if err != nil {
		uc.fileStatus.Failed(ctx, fileId, uptr.String(err.Error()))
		uc.saveStateByFileId(ctx, fileId)
		return nil, err
	}

	file, err = uc.file.FindByFileId(ctx, fileId, true)
	if err != nil {
		return nil, err
	}

	uc.saveStateByFileId(ctx, fileId)

	uc.enqueueDownloadTask(file.DownloadTask)

	return &dto.ScheduleDownloadResponse{
		FileId:       file.FileId,
		Status:       file.Status,
		YoutubeTitle: file.YoutubeTitle,
		Format:       file.Ext,
	}, nil
}
