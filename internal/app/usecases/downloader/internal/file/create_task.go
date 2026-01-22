package fileuc

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *File) CreateTask(ctx context.Context, file *ddownload.File, dlOptions *ddownload.DownloadOptions) error {
	task := &ddownload.DownloadTask{
		FileId:     file.FileId,
		YoutubeUrl: file.YoutubeUrl,
		Options:    dlOptions,
	}

	return uc.dlTask.Create(ctx, task)
}
