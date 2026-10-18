package downloader

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *downloader) ExecuteDownloadTask(
	ctx context.Context,
	workerId uint64,
	task *ddownload.DownloadTask,
) error {
	return uc.taskExecutors.Download.Execute(ctx, workerId, task)
}
