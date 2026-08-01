package downloader

import (
	"context"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *Downloader) ExecuteCreateMediaWatchEvent(
	ctx context.Context,
	workerID uint64,
	req *dto.CreateMediaWatchEventRequest,
) error {
	if req == nil || req.Event == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	err := uc.mediaWatch.ExecuteCreateMediaWatchEvent(ctx, workerID, req)
	if err != nil {
		return err
	}

	if req.EventType == dtypes.MediaWatchEventTypePause || req.EventType == dtypes.MediaWatchEventTypeEnded {
		uc.broadcastDownloadUpdateToAuth(ctx, req.Event.AuthCtx(), req.Event.DownloadID)
	}

	return nil
}
