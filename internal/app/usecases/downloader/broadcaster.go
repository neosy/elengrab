package downloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *YouTubeDownloader) Broadcaster() *broadcaster.Broadcaster {
	return uc.broadcaster
}

func (uc *YouTubeDownloader) broadcastFileAdd(file *ddownload.File) {
	if file == nil {
		return
	}

	var accessByUserID string
	if uc.historyMode != dtypes.HistoryModeGlobal && file.UserID != nil {
		accessByUserID = file.UserID.String()
	}

	resp := &dto.ScheduleDownloadResponse{
		URL:        file.MediaUrl,
		FileID:     file.FileID,
		Status:     file.Status,
		MediaTitle: file.MediaTitle,
		Format:     file.Ext,
	}

	if accessByUserID == "" {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeFileAdd, resp)
	} else {
		uc.broadcaster.BroadcastToUser(accessByUserID, dto.BroadcastEventTypeFileAdd, resp)
	}
}

func (uc *YouTubeDownloader) broadcastFileUpdate(
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

	var accessByUserID string
	if uc.historyMode != dtypes.HistoryModeGlobal && fileInfo.UserID != nil {
		accessByUserID = fileInfo.UserID.String()
	}

	if accessByUserID == "" {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeFileUpdate, fileInfo)
	} else {
		uc.broadcaster.BroadcastToUser(accessByUserID, dto.BroadcastEventTypeFileUpdate, fileInfo)
	}
}

func (uc *YouTubeDownloader) broadcastFileDelete(
	userID *uuid.UUID,
	fileID uuid.UUID,
) {
	var accessByUserID string
	if uc.historyMode != dtypes.HistoryModeGlobal && userID != nil {
		accessByUserID = userID.String()
	}

	if accessByUserID == "" {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeFileDelete, fileID)
	} else {
		uc.broadcaster.BroadcastToUser(accessByUserID, dto.BroadcastEventTypeFileDelete, fileID)
	}
}

func (uc *YouTubeDownloader) broadcastFileProgressUpdate(
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

	var accessByUserID string
	if uc.historyMode != dtypes.HistoryModeGlobal && fileInfo.UserID != nil {
		accessByUserID = fileInfo.UserID.String()
	}

	if accessByUserID == "" {
		uc.broadcaster.Broadcast(dto.BroadcastEventTypeProgressUpdate, resp)
	} else {
		uc.broadcaster.BroadcastToUser(accessByUserID, dto.BroadcastEventTypeProgressUpdate, resp)
	}
}
