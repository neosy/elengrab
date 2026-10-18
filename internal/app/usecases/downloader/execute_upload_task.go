package downloader

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *downloader) ExecuteUploadTask(
	ctx context.Context,
	workerId uint64,
	task *ddownload.MediaProcessingTask,
) error {
	return uc.taskExecutors.Upload.Execute(ctx, workerId, task)
}
