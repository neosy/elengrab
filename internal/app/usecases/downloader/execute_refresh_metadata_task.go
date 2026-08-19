package downloader

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *Downloader) ExecuteRefreshMetadataTask(
	ctx context.Context,
	workerID uint64,
	task *ddownload.RefreshMetadataTask,
) error {
	return uc.dlExecutor.ExecuteRefreshMetadataTask(ctx, workerID, task)
}
