package downloader

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *downloader) ExecuteRefreshMetadataTask(
	ctx context.Context,
	workerID uint64,
	task *ddownload.RefreshMetadataTask,
) error {
	return uc.taskExecutors.RefreshMetadata.Execute(ctx, workerID, task)
}
