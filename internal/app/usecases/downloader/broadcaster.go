package downloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *Downloader) Broadcaster() *broadcaster.Broadcaster {
	return uc.broadcaster
}

func (uc *Downloader) broadcastFileAdd(file *ddownload.File) {
	if file == nil {
		return
	}

	var accessByUserID uuid.UUID
	if uc.authz.RestrictFilesByUser(nil) && file.UserID != nil {
		accessByUserID = *file.UserID
	}

	resp := &dto.ScheduleDownloadResponse{
		URL:        file.MediaUrl,
		FileID:     file.FileID,
		Status:     file.Status,
		MediaTitle: file.MediaTitle,
		Format:     file.Ext,
	}

	if accessByUserID == uuid.Nil {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeFileAdd, resp)
	} else {
		uc.broadcaster.BroadcastToUser(accessByUserID, dto.BroadcastEventTypeFileAdd, resp)
	}
}

func (uc *Downloader) broadcastFileUpdate(
	ctx context.Context,
	fileID uuid.UUID,
) {
	fileInfo, err := uc.findActualFileInfo(ctx, nil, fileID)
	if err != nil {
		uc.logger.Error("Failed find file info", "error", err)
		return
	}

	if fileInfo == nil {
		return
	}

	var accessByUserID uuid.UUID
	if uc.authz.RestrictFilesByUser(nil) && fileInfo.UserID != nil {
		accessByUserID = *fileInfo.UserID
	}

	if accessByUserID == uuid.Nil {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeFileUpdate, fileInfo)
	} else {
		uc.broadcaster.BroadcastToUser(accessByUserID, dto.BroadcastEventTypeFileUpdate, fileInfo)
	}
}

func (uc *Downloader) broadcastFileDelete(
	userID *uuid.UUID,
	fileID uuid.UUID,
) {
	var accessByUserID uuid.UUID
	if uc.authz.RestrictFilesByUser(nil) && userID != nil {
		accessByUserID = *userID
	}

	if accessByUserID == uuid.Nil {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeFileDelete, fileID)
	} else {
		uc.broadcaster.BroadcastToUser(accessByUserID, dto.BroadcastEventTypeFileDelete, fileID)
	}
}

func (uc *Downloader) broadcastFileProgressUpdate(
	ctx context.Context,
	fileID uuid.UUID,
) {
	fileInfo, err := uc.findActualFileInfo(ctx, nil, fileID)
	if err != nil {
		uc.logger.Error("Failed find file info", "error", err)
		return
	}

	if fileInfo == nil {
		return
	}

	if fileInfo.WorkingStatus != dto.WorkingStatusDownloading {
		return
	}

	var percent int = 0
	if fileInfo.Progress != nil {
		percent = int(fileInfo.Progress.Percent())
	}

	resp := dto.FileProgressResponse{
		FileID:  fileID,
		Percent: percent,
	}

	var accessByUserID uuid.UUID
	if uc.authz.RestrictFilesByUser(nil) && fileInfo.UserID != nil {
		accessByUserID = *fileInfo.UserID
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
	userID uuid.UUID,
	module dto.BroadcastNotificationModule,
	notificationType dto.BroadcastNotificationType,
	message string,
) {
	notification := dto.BroadcastNotification{
		Module:  module,
		Type:    notificationType,
		Message: message,
	}
	uc.broadcaster.BroadcastToUser(userID, dto.BroadcastEventTypeNotification, notification)
}
