package downloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/app/usecases/types"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	eventkey "github.com/neosy/elengrab/internal/domain/types/event_key"
)

// newErrorBroadcaster returns a function that broadcasts the given error
// as an error notification for the caller and returns it unchanged.
func (uc *downloader) newErrorBroadcaster(
	authCtx dauth.AuthContext,
	broadcastNotificationModule types.BroadcastNotificationModule,
) func(err error) error {
	return func(err error) error {
		uc.broadcastNotification(
			authCtx.EventKey(),
			broadcastNotificationModule,
			types.BroadcastNotificationTypeError,
			err.Error(),
		)
		return err
	}
}

func (uc *downloader) Broadcaster() *broadcaster.Broadcaster {
	return uc.broadcaster
}

func (uc *downloader) broadcastMediaDownloadAdded(downloadInfo *dto.MediaDownloadInfo) {
	if downloadInfo == nil {
		return
	}

	if downloadInfo.UserID == nil || *downloadInfo.UserID == uuid.Nil {
		uc.broadcaster.Broadcast(types.BroadcastEventTypeDownloadAdd, downloadInfo)
		return
	}

	if downloadInfo.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			eventkey.NewEventKeyUserID(*downloadInfo.UserID),
			types.BroadcastEventTypeDownloadAdd,
			downloadInfo,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(*downloadInfo.UserID, types.BroadcastEventTypeDownloadAdd, downloadInfo)
}

func (uc *downloader) broadcastMediaDownloadUpdated(
	ctx context.Context,
	downloadID uuid.UUID,
) {
	downloadInfo, err := uc.findActualDownloadInfo(ctx, downloadID)
	if err != nil {
		uc.logger.Error("Failed find download info", "error", err)
		return
	}

	if downloadInfo == nil {
		return
	}

	if downloadInfo.UserID == nil || *downloadInfo.UserID == uuid.Nil {
		uc.broadcaster.Broadcast(types.BroadcastEventTypeDownloadUpdate, downloadInfo)
		return
	}

	if downloadInfo.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			eventkey.NewEventKeyUserID(*downloadInfo.UserID),
			types.BroadcastEventTypeDownloadUpdate,
			downloadInfo,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(*downloadInfo.UserID, types.BroadcastEventTypeDownloadUpdate, downloadInfo)
}

func (uc *downloader) broadcastWatchStatsUpdatedToAuth(
	ctx context.Context,
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
) {
	downloadInfo, err := uc.findActualDownloadInfo(ctx, downloadID, withAuth(authCtx))
	if err != nil {
		uc.logger.Error("Failed find download info", "error", err)
		return
	}

	if downloadInfo == nil {
		return
	}

	downloadChanged := &dto.MediaDownloadChanged{
		DownloadID: downloadID,
		Info:       downloadInfo,
	}

	downloadChanged.MarkWatchStatsChanged()

	uc.broadcaster.BroadcastToAuth(authCtx, types.BroadcastEventTypeDownloadPatch, downloadChanged)
}

func (uc *downloader) broadcastWatchPositionUpdatedToAuth(
	ctx context.Context,
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
) {
	downloadInfo, err := uc.findActualDownloadInfo(ctx, downloadID, withAuth(authCtx))
	if err != nil {
		uc.logger.Error("Failed find download info", "error", err)
		return
	}

	if downloadInfo == nil {
		return
	}

	downloadChanged := &dto.MediaDownloadChanged{
		DownloadID: downloadID,
		Info:       downloadInfo,
	}

	downloadChanged.MarkWatchPositionChanged()

	uc.broadcaster.BroadcastToAuth(authCtx, types.BroadcastEventTypeDownloadPatch, downloadChanged)
}

func (uc *downloader) broadcastMediaDownloadChanged(
	ctx context.Context,
	req dto.MediaDownloadChanged,
) {
	if req.Info == nil {
		downloadInfo, err := uc.findActualDownloadInfo(ctx, req.DownloadID)
		if err != nil {
			uc.logger.Error("Failed find download info", "error", err)
			return
		}
		req.Info = downloadInfo
	}

	if req.Info == nil {
		return
	}

	if req.Info.UserID == nil || *req.Info.UserID == uuid.Nil {
		uc.broadcaster.Broadcast(types.BroadcastEventTypeDownloadPatch, &req)
		return
	}

	if req.Info.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			eventkey.NewEventKeyUserID(*req.Info.UserID),
			types.BroadcastEventTypeDownloadPatch,
			&req,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(*req.Info.UserID, types.BroadcastEventTypeDownloadPatch, &req)
}

func (uc *downloader) broadcastMediaDownloadStartRefreshing(
	ctx context.Context,
	downloadID uuid.UUID,
) {
	downloadInfo, err := uc.findActualDownloadInfo(ctx, downloadID)
	if err != nil {
		uc.logger.Error("Failed find download info", "error", err)
		return
	}

	if downloadInfo == nil {
		return
	}

	if downloadInfo.UserID == nil || *downloadInfo.UserID == uuid.Nil {
		uc.broadcaster.Broadcast(types.BroadcastEventTypeDownloadStartRefreshing, downloadInfo)
		return
	}

	if downloadInfo.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			eventkey.NewEventKeyUserID(*downloadInfo.UserID),
			types.BroadcastEventTypeDownloadStartRefreshing,
			downloadInfo,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(*downloadInfo.UserID, types.BroadcastEventTypeDownloadStartRefreshing, downloadInfo)
}

func (uc *downloader) broadcastMediaDownloadDeleted(_ context.Context, download *ddownload.MediaDownload) {
	if download == nil {
		return
	}

	if download.UserID == nil || *download.UserID == uuid.Nil {
		uc.broadcaster.Broadcast(types.BroadcastEventTypeDownloadDelete, download.DownloadID)
		return
	}

	if download.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			eventkey.NewEventKeyUserID(*download.UserID),
			types.BroadcastEventTypeDownloadDelete,
			download.DownloadID,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(*download.UserID, types.BroadcastEventTypeDownloadDelete, download.DownloadID)

}

func (uc *downloader) broadcastDownloadProgressUpdated(
	ctx context.Context,
	downloadID uuid.UUID,
) {
	downloadInfo, err := uc.findActualDownloadInfo(ctx, downloadID)
	if err != nil {
		uc.logger.Error("Failed find download info", "error", err)
		return
	}

	if downloadInfo == nil {
		return
	}

	if downloadInfo.WorkingStatus != types.WorkingStatusDownloading {
		return
	}

	var percent int = 0
	if downloadInfo.Progress != nil {
		percent = int(downloadInfo.Progress.Percent())
	}

	resp := dto.MediaDownloadProgressResponse{
		DownloadID: downloadID,
		Percent:    percent,
	}

	if downloadInfo.UserID == nil || *downloadInfo.UserID == uuid.Nil {
		uc.broadcaster.Broadcast(types.BroadcastEventTypeProgressUpdate, resp)
		return
	}

	if downloadInfo.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			eventkey.NewEventKeyUserID(*downloadInfo.UserID),
			types.BroadcastEventTypeProgressUpdate,
			resp,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(*downloadInfo.UserID, types.BroadcastEventTypeProgressUpdate, resp)
}

func (uc *downloader) broadcastSystemInfoUpdated() {
	uc.broadcaster.Broadcast(types.BroadcastEventTypeSystemInfoUpdate, uc.SystemInfo())
}

func (uc *downloader) broadcastNotification(
	eventKey eventkey.EventKey,
	module types.BroadcastNotificationModule,
	notificationType types.BroadcastNotificationType,
	message string,
) {
	notification := dto.BroadcastNotification{
		Module:  module,
		Type:    notificationType,
		Message: message,
	}
	uc.broadcaster.BroadcastByKey(eventKey, types.BroadcastEventTypeNotification, notification)
}

func (uc *downloader) NotifyDownloadUpdated(
	ctx context.Context,
	downloadID uuid.UUID,
) {
	uc.broadcastMediaDownloadUpdated(ctx, downloadID)
}

func (uc *downloader) NotifyDownloadChanged(
	ctx context.Context,
	req dto.MediaDownloadChanged,
) {
	uc.broadcastMediaDownloadChanged(ctx, req)
}
