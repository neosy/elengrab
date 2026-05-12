package downloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *Downloader) Broadcaster() *broadcaster.Broadcaster {
	return uc.broadcaster
}

func (uc *Downloader) broadcastDownloadAdd(download *ddownload.MediaDownload) {
	if download == nil {
		return
	}

	var accessByUserID uuid.UUID
	if uc.authz.RestrictDownloadsByUser(nil) && download.UserID != nil {
		accessByUserID = *download.UserID
	}

	resp := &dto.ScheduleDownloadResponse{
		URL:        download.MediaURL,
		DownloadID: download.DownloadID,
		Status:     download.Status,
		MediaTitle: download.MediaTitle,
		Format:     download.Ext,
	}

	if accessByUserID == uuid.Nil {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeDownloadAdd, resp)
	} else {
		uc.broadcaster.BroadcastToUser(accessByUserID, dto.BroadcastEventTypeDownloadAdd, resp)
	}
}

func (uc *Downloader) broadcastDownloadUpdate(
	ctx context.Context,
	downloadID uuid.UUID,
) {
	downloadInfo, err := uc.findActualDownloadInfo(ctx, nil, downloadID)
	if err != nil {
		uc.logger.Error("Failed find download info", "error", err)
		return
	}

	if downloadInfo == nil {
		return
	}

	var accessByUserID uuid.UUID
	if uc.authz.RestrictDownloadsByUser(nil) && downloadInfo.UserID != nil {
		accessByUserID = *downloadInfo.UserID
	}

	if accessByUserID == uuid.Nil {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeDownloadUpdate, downloadInfo)
	} else {
		uc.broadcaster.BroadcastToUser(accessByUserID, dto.BroadcastEventTypeDownloadUpdate, downloadInfo)
	}
}

func (uc *Downloader) broadcastDownloadDelete(
	userID *uuid.UUID,
	downloadID uuid.UUID,
) {
	var accessByUserID uuid.UUID
	if uc.authz.RestrictDownloadsByUser(nil) && userID != nil {
		accessByUserID = *userID
	}

	if accessByUserID == uuid.Nil {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeDownloadDelete, downloadID)
	} else {
		uc.broadcaster.BroadcastToUser(accessByUserID, dto.BroadcastEventTypeDownloadDelete, downloadID)
	}
}

func (uc *Downloader) broadcastDownloadProgressUpdate(
	ctx context.Context,
	downloadID uuid.UUID,
) {
	downloadInfo, err := uc.findActualDownloadInfo(ctx, nil, downloadID)
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

	var accessByUserID uuid.UUID
	if uc.authz.RestrictDownloadsByUser(nil) && downloadInfo.UserID != nil {
		accessByUserID = *downloadInfo.UserID
	}

	if accessByUserID == uuid.Nil {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeProgressUpdate, resp)
	} else {
		uc.broadcaster.BroadcastToUser(accessByUserID, dto.BroadcastEventTypeProgressUpdate, resp)
	}
}

func (uc *Downloader) broadcastSystemInfoUpdate() {
	uc.broadcaster.Broadcast(dto.BroadcastEventTypeSystemInfoUpdate, uc.SystemInfo())
}

func (uc *Downloader) broadcastNotification(
	eventKey dtypes.EventKey,
	module dto.BroadcastNotificationModule,
	notificationType dto.BroadcastNotificationType,
	message string,
) {
	notification := dto.BroadcastNotification{
		Module:  module,
		Type:    notificationType,
		Message: message,
	}
	uc.broadcaster.BroadcastTo(eventKey, dto.BroadcastEventTypeNotification, notification)
}
