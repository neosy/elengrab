package downloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	eventkey "github.com/neosy/elengrab/internal/domain/types/event_key"
)

func (uc *Downloader) Broadcaster() *broadcaster.Broadcaster {
	return uc.broadcaster
}

func (uc *Downloader) broadcastDownloadAdd(download *ddownload.MediaDownload) {
	if download == nil {
		return
	}

	resp := &dto.ScheduleDownloadResponse{
		URL:        download.MediaURL,
		DownloadID: download.DownloadID,
		Status:     download.Status,
		MediaTitle: download.MediaTitle,
		Format:     download.Ext,
	}

	if download.UserID == nil || *download.UserID == uuid.Nil {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeDownloadAdd, resp)
		return
	}

	if download.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			eventkey.NewEventKeyUserID(*download.UserID),
			dto.BroadcastEventTypeDownloadAdd,
			resp,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(*download.UserID, dto.BroadcastEventTypeDownloadAdd, resp)
}

func (uc *Downloader) broadcastDownloadUpdate(
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
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeDownloadUpdate, downloadInfo)
		return
	}

	if downloadInfo.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			eventkey.NewEventKeyUserID(*downloadInfo.UserID),
			dto.BroadcastEventTypeDownloadUpdate,
			downloadInfo,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(*downloadInfo.UserID, dto.BroadcastEventTypeDownloadUpdate, downloadInfo)
}

func (uc *Downloader) broadcastDownloadStartRefreshing(
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
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeDownloadStartRefreshing, downloadInfo)
		return
	}

	if downloadInfo.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			eventkey.NewEventKeyUserID(*downloadInfo.UserID),
			dto.BroadcastEventTypeDownloadStartRefreshing,
			downloadInfo,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(*downloadInfo.UserID, dto.BroadcastEventTypeDownloadStartRefreshing, downloadInfo)
}

func (uc *Downloader) broadcastDownloadDelete(
	ctx context.Context,
	download *ddownload.MediaDownload,
) {
	if download == nil {
		return
	}

	if download.UserID == nil || *download.UserID == uuid.Nil {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeDownloadDelete, download.DownloadID)
		return
	}

	if download.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			eventkey.NewEventKeyUserID(*download.UserID),
			dto.BroadcastEventTypeDownloadDelete,
			download.DownloadID,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(*download.UserID, dto.BroadcastEventTypeDownloadDelete, download.DownloadID)

}

func (uc *Downloader) broadcastDownloadProgressUpdate(
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

	if downloadInfo.WorkingStatus != dto.WorkingStatusDownloading {
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
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeProgressUpdate, resp)
		return
	}

	if downloadInfo.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			eventkey.NewEventKeyUserID(*downloadInfo.UserID),
			dto.BroadcastEventTypeProgressUpdate,
			resp,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(*downloadInfo.UserID, dto.BroadcastEventTypeProgressUpdate, resp)
}

func (uc *Downloader) broadcastSystemInfoUpdate() {
	uc.broadcaster.Broadcast(dto.BroadcastEventTypeSystemInfoUpdate, uc.SystemInfo())
}

func (uc *Downloader) broadcastNotification(
	eventKey eventkey.EventKey,
	module dto.BroadcastNotificationModule,
	notificationType dto.BroadcastNotificationType,
	message string,
) {
	notification := dto.BroadcastNotification{
		Module:  module,
		Type:    notificationType,
		Message: message,
	}
	uc.broadcaster.BroadcastByKey(eventKey, dto.BroadcastEventTypeNotification, notification)
}

func (uc *Downloader) NotifyDownloadUpdated(
	ctx context.Context,
	downloadID uuid.UUID,
) {
	uc.broadcastDownloadUpdate(ctx, downloadID)
}
