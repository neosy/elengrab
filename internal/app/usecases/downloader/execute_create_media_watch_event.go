package downloader

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

func (uc *Downloader) ExecuteCreateMediaWatchEvent(
	ctx context.Context,
	workerID uint64,
	req *dto.CreateMediaWatchEventRequest,
) error {
	return uc.mediaWatch.ExecuteCreateMediaWatchEvent(ctx, workerID, req)
}
