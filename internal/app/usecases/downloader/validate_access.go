package downloader

import (
	"slices"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/authz"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	iconfig "github.com/neosy/elengrab/internal/config"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/exceptions"
)

func (uc *downloader) CanCreateMediaDownload(authCtx dauth.AuthContext) bool {
	if iconfig.DemoMode() {
		return false
	}

	if uc.authz.HasCreateAccess(authCtx.RoleIDs) {
		return true
	}

	if iconfig.AppMode() == dtypes.AppModeGuest {
		return true
	}

	return false
}

// validateWriteOperationAccess validates whether the user is allowed to perform write operations.
func (uc *downloader) validateWriteOperationAccess(authCtx dauth.AuthContext) error {
	if uc.demoMode {
		uc.broadcastNotification(
			authCtx.EventKey(),
			dto.BroadcastNotificationModuleResultRow,
			dto.BroadcastNotificationTypeError,
			"Operation not allowed in demo mode",
		)
		return exceptions.DEMO_MODE_RESTRICTION.NewErrorx()
	}

	if !uc.authz.HasFullAccess(authCtx.RoleIDs) && authz.IsAnonymous(authCtx.RoleIDs) {
		uc.broadcastNotification(
			authCtx.EventKey(),
			dto.BroadcastNotificationModuleResultRow,
			dto.BroadcastNotificationTypeError,
			"You must be authenticated to perform this action",
		)
		return ierrors.ErrUnauthorized
	}

	return nil
}

func (uc *downloader) validateDownloadEditAccess(authCtx dauth.AuthContext, download *ddownload.MediaDownload) error {
	if iconfig.DemoMode() {
		return ierrors.ErrDemoModeAccessDenied
	}

	if !slices.Contains(dtypes.MediaDownloadEditableStatuses(), download.Status) {
		return ierrors.ErrAccessDenied
	}
	if uc.authz.HasFullAccess(authCtx.RoleIDs) {
		return nil
	}

	if download.UserID == nil {
		return nil
	}

	if *download.UserID == authCtx.UserID {
		return nil
	}

	return ierrors.ErrAccessDenied
}

func (uc *downloader) validateDownloadDeleteAccess(authCtx dauth.AuthContext, download *ddownload.MediaDownload) error {
	if iconfig.DemoMode() {
		return ierrors.ErrDemoModeAccessDenied
	}

	return uc.validateDownloadEditAccess(authCtx, download)
}

// HasWriteOperationAccess reports whether the user is allowed to perform operations other than read-only access.
func (uc *downloader) HasWriteOperationAccess(authCtx dauth.AuthContext) bool {
	return uc.validateWriteOperationAccess(authCtx) == nil
}

// HasDownloadEditAccess reports whether the user is allowed to edit the specified download.
func (uc *downloader) HasDownloadEditAccess(authCtx dauth.AuthContext, download *ddownload.MediaDownload) bool {
	return uc.validateDownloadEditAccess(authCtx, download) == nil
}

// HasDownloadDeleteAccess reports whether the user is allowed to delete the specified download.
func (uc *downloader) HasDownloadDeleteAccess(authCtx dauth.AuthContext, download *ddownload.MediaDownload) bool {
	return uc.validateDownloadDeleteAccess(authCtx, download) == nil
}
