package downloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/valyala/fasthttp"
)

// RepeatDownload repeats the download process for a specific file.
func (uc *YouTubeDownloader) RepeatDownload(
	ctx context.Context,
	userID uuid.UUID,
	fileID uuid.UUID,
) (*dto.GetFileInfoResponse, error) {
	var accessByUserID *uuid.UUID
	if uc.historyMode != dtypes.HistoryModeGlobal {
		accessByUserID = &userID
	}

	if uc.demoMode {
		uc.broadcastNotification(
			userID,
			dto.BroadcastNotificationModuleResultRow,
			dto.BroadcastNotificationTypeError,
			"Operation not allowed in demo mode",
		)
		return nil, errorx.NewHTTP("operation not allowed in demo mode", fasthttp.StatusForbidden)
	}

	err := uc.file.Tx(
		ctx,
		func(ctx context.Context) error {
			file, err := uc.file.GetByFileID(ctx, accessByUserID, fileID)
			if err != nil {
				return err
			}

			err = uc.fileStatus.New(ctx, file.FileID)
			if err != nil {
				uc.logger.Error("Failed to set file status to new", "error", err)
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	file, err := uc.file.GetByFileID(ctx, accessByUserID, fileID)
	if err != nil {
		return nil, err
	}

	err = uc.addFileToQueueDownload(ctx, fileID, file.DownloadTask.TaskID)
	if err != nil {
		uc.logger.Error("Failed add to queue", "error", err)
		return nil, err
	}

	file, err = uc.file.GetByFileID(ctx, accessByUserID, fileID)
	if err != nil {
		return nil, err
	}

	uc.dlStateCache.SaveByFile(ctx, file)
	uc.broadcastFileUpdate(ctx, fileID)

	return uc.findActualFileInfoByFile(ctx, file)
}
