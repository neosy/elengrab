package mediadownload

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaDownload) CreateTask(ctx context.Context, download *ddownload.MediaDownload, dlOptions *ddownload.DownloadOptions) error {
	task := &ddownload.DownloadTask{
		DownloadID: download.DownloadID,
		MediaUrl:   download.MediaURL,
		Options:    dlOptions,
	}

	return uc.dlTask.Create(ctx, task)
}
