package downloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

// RetryDownload repeats the download process for a specific download.
func (uc *Downloader) RetryDownload(
	ctx context.Context,
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
) (*dto.GetMediaDownloadInfoResponse, error) {
	err := uc.validateWriteOperation(authCtx)
	if err != nil {
		return nil, err
	}

	resetStatusToNew := func(ctx context.Context) error {
		download, err := uc.download.GetByDownloadIDNoCache(ctx, downloadID)
		if err != nil {
			return err
		}

		err = uc.validateDownloadWriteAccess(authCtx, download)
		if err != nil {
			return err
		}

		err = uc.downloadStatus.New(ctx, download.DownloadID)
		if err != nil {
			uc.logger.Error("Failed to set download status to new", "error", err)
			return err
		}
		return nil
	}

	err = uc.download.Tx(ctx, resetStatusToNew)
	if err != nil {
		return nil, err
	}

	download, err := uc.download.GetByDownloadIDNoCache(ctx, downloadID)
	if err != nil {
		return nil, err
	}

	err = uc.addDownloadToQueueDownload(ctx, downloadID, download.DownloadTask.TaskID)
	if err != nil {
		uc.logger.Error("Failed add to queue", "error", err)
		return nil, err
	}

	download, err = uc.download.GetByDownloadIDNoCache(ctx, downloadID)
	if err != nil {
		return nil, err
	}

	uc.dlStateCache.SaveByDownload(ctx, download)
	uc.broadcastDownloadUpdate(ctx, downloadID)

	return uc.findActualDownloadInfoByDownload(ctx, download)
}
