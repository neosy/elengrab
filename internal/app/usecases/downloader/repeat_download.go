package downloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *YouTubeDownloader) RepeatDownload(
	ctx context.Context,
	userID uuid.UUID,
	fileID uuid.UUID,
) (*dto.GetFileInfoResponse, error) {
	var accessByUserID *uuid.UUID
	if uc.historyMode != dtypes.HistoryModeGlobal {
		accessByUserID = &userID
	}

	err := uc.file.Tx(
		ctx,
		func(ctx context.Context) error {
			file, err := uc.file.GetByFileId(ctx, accessByUserID, fileID)
			if err != nil {
				return err
			}

			err = uc.fileStatus.New(ctx, file.FileId)
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

	file, err := uc.file.GetByFileId(ctx, accessByUserID, fileID)
	if err != nil {
		return nil, err
	}

	err = uc.addFileToQueueDownload(ctx, fileID, file.DownloadTask.TaskId)
	if err != nil {
		uc.logger.Error("Failed add to queue", "error", err)
		return nil, err
	}

	file, err = uc.file.GetByFileId(ctx, accessByUserID, fileID)
	if err != nil {
		return nil, err
	}

	uc.dlStateCache.SaveByFile(ctx, file)

	return uc.findStateAndFileInfo(ctx, nil, nil, file, false)
}
