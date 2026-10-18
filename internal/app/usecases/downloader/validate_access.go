package downloader

import (
	"slices"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/authz"
	"github.com/neosy/elengrab/internal/app/usecases/types"
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

// validateWriteAccess returns an error if the caller has no access to write operations
// in the current app mode.
func (uc *downloader) validateWriteAccess(
	authCtx dauth.AuthContext,
	broadcastNotificationModule types.BroadcastNotificationModule,
) error {
	broadcastTextError := func(text string) {
		if broadcastNotificationModule == types.BroadcastNotificationModuleNone {
			return
		}

		uc.broadcastNotification(
			authCtx.EventKey(),
			broadcastNotificationModule,
			types.BroadcastNotificationTypeError,
			text,
		)
	}

	if uc.demoMode {
		broadcastTextError("Operation not allowed in demo mode")
		return exceptions.DEMO_MODE_RESTRICTION.NewErrorx()
	}

	if uc.appMode.RequiresIdentityForWrite() && authCtx.IsAnonymous() {
		broadcastTextError("You must be authenticated to perform this action")
		return exceptions.UNAUTHORIZED.NewErrorx()
	}

	return nil
}

func (uc *downloader) validateDownloadEditableStatus(download *ddownload.MediaDownload) error {
	if !slices.Contains(dtypes.MediaDownloadEditableStatuses(), download.Status) {
		return ierrors.ErrMediaDownloadNotEditable
	}

	return nil
}

func (uc *downloader) validateDownloadViewAccess(
	authCtx dauth.AuthContext,
	download *ddownload.MediaDownload,
) error {
	if !uc.authz.HasMediaDownloadViewAccess(authCtx, download) {
		return ierrors.ErrAccessDenied
	}

	return nil
}

func (uc *downloader) validateDownloadEditAccess(
	authCtx dauth.AuthContext,
	download *ddownload.MediaDownload,
) error {
	if iconfig.DemoMode() {
		return ierrors.ErrDemoModeAccessDenied
	}

	if err := uc.validateDownloadEditableStatus(download); err != nil {
		return err
	}

	if err := uc.validateDownloadViewAccess(authCtx, download); err != nil {
		return err
	}

	if uc.authz.HasFullAccess(authCtx.RoleIDs) {
		return nil
	}

	if download.UserID == nil || authz.IsAnonymous(authCtx.RoleIDs) {
		return ierrors.ErrAccessDenied
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

	if err := uc.validateDownloadViewAccess(authCtx, download); err != nil {
		return err
	}

	if uc.authz.HasFullAccess(authCtx.RoleIDs) {
		return nil
	}

	if download.UserID == nil || authCtx.IsAnonymous() {
		return ierrors.ErrAccessDenied
	}

	if *download.UserID == authCtx.UserID {
		return nil
	}

	return ierrors.ErrAccessDenied
}

func (uc *downloader) validateDownloadRetryAccess(authCtx dauth.AuthContext, download *ddownload.MediaDownload) error {
	if iconfig.DemoMode() {
		return ierrors.ErrDemoModeAccessDenied
	}

	if download.Status != dtypes.MediaDownloadStatusFailed {
		return ierrors.ErrMediaDownloadNotRetryable
	}

	return uc.validateDownloadDeleteAccess(authCtx, download)
}

// HasWriteAccess reports whether the user is allowed to perform operations other than read-only access.
func (uc *downloader) HasWriteAccess(authCtx dauth.AuthContext) bool {
	return uc.validateWriteAccess(authCtx, types.BroadcastNotificationModuleNone) == nil
}

// HasDownloadEditAccess reports whether the user is allowed to edit the specified download.
func (uc *downloader) HasDownloadEditAccess(authCtx dauth.AuthContext, download *ddownload.MediaDownload) bool {
	return uc.validateDownloadEditAccess(authCtx, download) == nil
}

// HasDownloadDeleteAccess reports whether the user is allowed to delete the specified download.
func (uc *downloader) HasDownloadDeleteAccess(authCtx dauth.AuthContext, download *ddownload.MediaDownload) bool {
	return uc.validateDownloadDeleteAccess(authCtx, download) == nil
}
