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

func (uc *Downloader) CanAddMediaDownload(authCtx dauth.AuthContext) bool {
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

func (uc *Downloader) HasWriteOperation(authCtx dauth.AuthContext) bool {
	if iconfig.DemoMode() {
		return false
	}

	if !uc.authz.HasWriteAllAccess(authCtx.RoleIDs) && authz.IsAnonymous(authCtx.RoleIDs) {
		return false
	}

	return true
}

func (uc *Downloader) validateWriteOperation(authCtx dauth.AuthContext) error {
	if uc.demoMode {
		uc.broadcastNotification(
			authCtx.EventKey(),
			dto.BroadcastNotificationModuleResultRow,
			dto.BroadcastNotificationTypeError,
			"Operation not allowed in demo mode",
		)
		return exceptions.DEMO_MODE_RESTRICTION.NewErrorx()
	}

	if !uc.authz.HasWriteAllAccess(authCtx.RoleIDs) && authz.IsAnonymous(authCtx.RoleIDs) {
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

func (uc *Downloader) validateDownloadWriteAccess(authCtx dauth.AuthContext, download *ddownload.MediaDownload) error {
	if uc.authz.HasWriteAllAccess(authCtx.RoleIDs) {
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

func (uc *Downloader) validateDownloadEditAccess(authCtx dauth.AuthContext, download *ddownload.MediaDownload) error {
	if !slices.Contains(dtypes.MediaDownloadEditableStatuses(), download.Status) {
		return ierrors.ErrAccessDenied
	}
	return uc.validateDownloadWriteAccess(authCtx, download)
}

func (uc *Downloader) validateDownloadDeleteAccess(authCtx dauth.AuthContext, download *ddownload.MediaDownload) error {
	return uc.validateDownloadWriteAccess(authCtx, download)
}
