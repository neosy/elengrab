package mediadownload

import (
	"context"
	"errors"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaDownload) CreateDownloadTask(
	ctx context.Context,
	download *ddownload.MediaDownload,
	dlOptions *ddownload.DownloadOptions,
) error {
	if download.MediaURL == nil {
		return errors.New("media url is required to create a download task")
	}

	task := &ddownload.DownloadTask{
		DownloadID: download.DownloadID,
		MediaUrl:   *download.MediaURL,
		Options:    dlOptions,
	}

	return uc.dlTask.Create(ctx, task)
}

func (uc *MediaDownload) CreateUploadTask(
	ctx context.Context,
	download *ddownload.MediaDownload,
	sourcePath string,
) error {
	task := &ddownload.MediaProcessingTask{
		DownloadID: download.DownloadID,
		SourcePath: sourcePath,
	}

	return uc.processingTask.Create(ctx, task)
}
