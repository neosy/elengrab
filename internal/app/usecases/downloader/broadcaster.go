package downloader

import (
	"context"

	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/app/usecases/types"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	eventkey "github.com/neosy/elengrab/internal/domain/types/event_key"
	dsse "github.com/neosy/elengrab/internal/domain/types/sse"
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

func (uc *downloader) broadcastMediaDownloadAdded(
	ctx context.Context,
	authCtx dauth.AuthContext,
	downloadInfo *dto.MediaDownloadInfo,
) {
	if downloadInfo == nil {
		return
	}

	if !downloadInfo.Status.IsReady() {
		if authCtx.HasSession() {
			uc.broadcaster.BroadcastToAuthorized(authCtx, types.BroadcastEventTypeDownloadAdd, downloadInfo)
			uc.broadcaster.BroadcastToViewAllAccess(ctx, authCtx.EventKey(), types.BroadcastEventTypeDownloadAdd, downloadInfo)
			return
		}

		if dauth.IsNotRegisteredUserID(downloadInfo.UserID) {
			return
		}
	}

	if downloadInfo.Status.IsReady() {
		if downloadInfo.Visibility == dtypes.MediaVisibilityPublic {
			uc.broadcaster.BroadcastPublic(
				ctx,
				eventkey.NewEventKeyUserID(downloadInfo.UserID),
				types.BroadcastEventTypeDownloadAdd,
				downloadInfo,
			)
			return
		}

		if dauth.IsNotRegisteredUserID(downloadInfo.UserID) {
			uc.broadcaster.Broadcast(types.BroadcastEventTypeDownloadAdd, downloadInfo)
			return
		}
	}

	uc.broadcaster.BroadcastToUsersWithAccess(ctx, downloadInfo.UserID, types.BroadcastEventTypeDownloadAdd, downloadInfo)
}

func (uc *downloader) broadcastMediaDownloadAddedForUser(
	ctx context.Context,
	userID uuid.UUID,
	download *ddownload.MediaDownload,
) {
	authCtx := dauth.AuthContext{
		UserID: userID,
	}

	downloadInfo, err := uc.resolveActualDownloadInfoByDownload(ctx, download)
	if err != nil {
		return
	}

	uc.broadcastMediaDownloadAdded(ctx, authCtx, downloadInfo)
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

	if downloadInfo.UserID == uuid.Nil() {
		uc.broadcaster.Broadcast(types.BroadcastEventTypeDownloadUpdate, downloadInfo)
		return
	}

	if downloadInfo.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			ctx,
			eventkey.NewEventKeyUserID(downloadInfo.UserID),
			types.BroadcastEventTypeDownloadUpdate,
			downloadInfo,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(ctx, downloadInfo.UserID, types.BroadcastEventTypeDownloadUpdate, downloadInfo)
}

func (uc *downloader) broadcastWatchStatsUpdated(
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

	downloadChanged := dto.MediaDownloadChanged{
		DownloadID: downloadID,
		Info:       downloadInfo,
	}

	downloadChanged.MarkUserWatchStatsChanged()
	uc.broadcaster.BroadcastToAuthorized(authCtx, types.BroadcastEventTypeDownloadPatch, downloadChanged)

	downloadChanged.ResetMark()
	downloadChanged.MarkPublicWatchStatsChanged()

	if downloadInfo.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			ctx,
			eventkey.NewEventKeyUserID(downloadInfo.UserID),
			types.BroadcastEventTypeDownloadPatch,
			downloadChanged,
		)
		return
	}

	uc.broadcaster.BroadcastToViewAllAccess(ctx, authCtx.EventKey(), types.BroadcastEventTypeDownloadPatch, downloadChanged)
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

	uc.broadcaster.BroadcastToAuthorized(authCtx, types.BroadcastEventTypeDownloadPatch, downloadChanged)
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

	if req.Info.UserID == uuid.Nil() {
		uc.broadcaster.Broadcast(types.BroadcastEventTypeDownloadPatch, &req)
		return
	}

	if req.Info.Status.IsReady() && req.Info.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			ctx,
			eventkey.NewEventKeyUserID(req.Info.UserID),
			types.BroadcastEventTypeDownloadPatch,
			&req,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(ctx, req.Info.UserID, types.BroadcastEventTypeDownloadPatch, &req)
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

	if downloadInfo.UserID == uuid.Nil() {
		uc.broadcaster.Broadcast(types.BroadcastEventTypeDownloadStartRefreshing, downloadInfo)
		return
	}

	if downloadInfo.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			ctx,
			eventkey.NewEventKeyUserID(downloadInfo.UserID),
			types.BroadcastEventTypeDownloadStartRefreshing,
			downloadInfo,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(ctx, downloadInfo.UserID, types.BroadcastEventTypeDownloadStartRefreshing, downloadInfo)
}

func (uc *downloader) broadcastMediaDownloadDeleted(ctx context.Context, download *ddownload.MediaDownload) {
	if download == nil {
		return
	}

	if download.UserID == nil || *download.UserID == uuid.Nil() {
		uc.broadcaster.Broadcast(types.BroadcastEventTypeDownloadDelete, download.DownloadID)
		return
	}

	if download.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			ctx,
			eventkey.NewEventKeyUserID(*download.UserID),
			types.BroadcastEventTypeDownloadDelete,
			download.DownloadID,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(ctx, *download.UserID, types.BroadcastEventTypeDownloadDelete, download.DownloadID)

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

	if downloadInfo.UserID == uuid.Nil() {
		uc.broadcaster.Broadcast(types.BroadcastEventTypeProgressUpdate, resp)
		return
	}

	if downloadInfo.Visibility == dtypes.MediaVisibilityPublic {
		uc.broadcaster.BroadcastPublic(
			ctx,
			eventkey.NewEventKeyUserID(downloadInfo.UserID),
			types.BroadcastEventTypeProgressUpdate,
			resp,
		)
		return
	}

	uc.broadcaster.BroadcastToUsersWithAccess(ctx, downloadInfo.UserID, types.BroadcastEventTypeProgressUpdate, resp)
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

// GetDownloadsInProgressEvent returns the SSE event for in-progress downloads.
func (uc *downloader) GetDownloadsInProgressEvent(
	ctx context.Context,
	authCtx dauth.AuthContext,
) (dsse.BroadcastEvent, error) {
	eventType := types.BroadcastEventTypeDownloadsInProgress

	if !authCtx.HasSession() {
		return dsse.EmptyBroadcastEvent(), nil
	}

	downloadsCount, err := uc.GetIncompleteDownloadsCount(ctx, authCtx.UserID)
	if err != nil {
		return dsse.EmptyBroadcastEvent(), err
	}

	if !uc.authz.HasViewAllAccess(authCtx.RoleIDs) {
		downloadsCount.ClearTotal()
	}

	return dsse.NewBroadcastEvent(eventType, downloadsCount), nil
}

// broadcastDownloadsInProgressToUser broadcasts in-progress downloads only to the specified user.
func (uc *downloader) broadcastDownloadsInProgressToUser(
	ctx context.Context,
	authCtx dauth.AuthContext,
) {
	event, _ := uc.GetDownloadsInProgressEvent(ctx, authCtx)
	if event.IsEmpty() {
		return
	}

	uc.broadcaster.BroadcastToAuthorized(authCtx, event.Type, event.Data)
}

// broadcastDownloadsInProgress broadcasts in-progress downloads to the specified user and users with extended access.
func (uc *downloader) broadcastDownloadsInProgress(ctx context.Context, authCtx dauth.AuthContext) {
	uc.broadcastDownloadsInProgressToUser(ctx, authCtx)

	uc.broadcaster.BroadcastToViewAllAccessWithDataBuilder(
		ctx, authCtx.EventKey(), types.BroadcastEventTypeDownloadsInProgress,
		func(userID uuid.UUID, roleIDs dtypes.UserRoleIDs) (any, bool) {
			data, err := uc.GetIncompleteDownloadsCount(ctx, userID)
			if err != nil {
				return nil, false
			}

			return data, true
		},
	)
}

// broadcastDownloadsInProgressForUser broadcasts in-progress downloads to the specified user and users with extended access.
func (uc *downloader) broadcastDownloadsInProgressForUser(ctx context.Context, userID uuid.UUID) {
	authCtx := dauth.AuthContext{
		UserID: userID,
	}

	if userID != dauth.AnonymousUserID() {
		user, _ := uc.authSrv.FindByUserID(ctx, userID)
		if user != nil {
			authCtx.RoleIDs, _ = user.ParseRoleIDs()
		}
	}

	uc.broadcastDownloadsInProgress(ctx, authCtx)
}

// NotifyDownloadUpdated notifies users about an updated media download.
func (uc *downloader) NotifyDownloadUpdated(
	ctx context.Context,
	downloadID uuid.UUID,
) {
	uc.broadcastMediaDownloadUpdated(ctx, downloadID)
}

// NotifyDownloadChanged notifies relevant users about changes to a media download.
func (uc *downloader) NotifyDownloadChanged(
	ctx context.Context,
	req dto.MediaDownloadChanged,
) {
	uc.broadcastMediaDownloadChanged(ctx, req)
}

// NotifyUserDownloadsInProgress notifies the specified user about in-progress downloads.
func (uc *downloader) NotifyUserDownloadsInProgress(
	ctx context.Context,
	authCtx dauth.AuthContext,
) {
	uc.broadcastDownloadsInProgressToUser(ctx, authCtx)
}
