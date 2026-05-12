package downloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/exceptions"
)

// RepeatDownload repeats the download process for a specific download.
func (uc *Downloader) RepeatDownload(
	ctx context.Context,
	userCtx dauth.UserContext,
	downloadID uuid.UUID,
) (*dto.GetMediaDownloadInfoResponse, error) {
	var accessByUserID *uuid.UUID
	if uc.authz.RestrictDownloadsByUser(userCtx.Roles) {
		accessByUserID = &userCtx.UserID
	}

	if uc.demoMode {
		uc.broadcastNotification(
			userCtx.EventKey(),
			dto.BroadcastNotificationModuleResultRow,
			dto.BroadcastNotificationTypeError,
			"Operation not allowed in demo mode",
		)
		return nil, exceptions.DEMO_MODE_RESTRICTION.NewErrorx()
	}

	err := uc.download.Tx(
		ctx,
		func(ctx context.Context) error {
			download, err := uc.download.GetByDownloadID(ctx, accessByUserID, downloadID)
			if err != nil {
				return err
			}

			err = uc.downloadStatus.New(ctx, download.DownloadID)
			if err != nil {
				uc.logger.Error("Failed to set download status to new", "error", err)
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	download, err := uc.download.GetByDownloadID(ctx, accessByUserID, downloadID)
	if err != nil {
		return nil, err
	}

	err = uc.addDownloadToQueueDownload(ctx, downloadID, download.DownloadTask.TaskID)
	if err != nil {
		uc.logger.Error("Failed add to queue", "error", err)
		return nil, err
	}

	download, err = uc.download.GetByDownloadID(ctx, accessByUserID, downloadID)
	if err != nil {
		return nil, err
	}

	uc.dlStateCache.SaveByDownload(ctx, download)
	uc.broadcastDownloadUpdate(ctx, downloadID)

	return uc.findActualDownloadInfoByDownload(ctx, download)
}
