package mediadownload

import (
	"context"
	"errors"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaDownload) CreateTask(
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
