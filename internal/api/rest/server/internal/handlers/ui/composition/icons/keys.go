package icons

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const (
	UserAvatarAdminIconNameKey  = "UserAvatarAdminIconName"
	UserAvatarUserIconNameKey   = "UserAvatarUserIconName"
	UserAvatarGuestIconNameKey  = "UserAvatarGuestIconName"
	UserAvatarAnonymIconNameKey = "UserAvatarAnonymIconName"

	DownloadIconNameKey           = "DownloadIconName"
	DownloadFailedIconNameKey     = "DownloadFailedIconName"
	DownloadPendingIconNameKey    = "DownloadPendingIconName"
	DownloadDeleteIconNameKey     = "DownloadDeleteIconName"
	MediaDefaultIconNameKey       = "MediaDefaultIconName"
	DownloadRepeatIconNameKey     = "DownloadRepeatIconName"
	DownloadSourceLinkIconNameKey = "DownloadSourceLinkIconName"
)

var (
	iconFileNames = map[string]any{
		UserAvatarAdminIconNameKey:  "user-admin-2.svg",
		UserAvatarUserIconNameKey:   "user-default.svg",
		UserAvatarGuestIconNameKey:  "user-guest.svg",
		UserAvatarAnonymIconNameKey: "user-anonymous-2.svg",

		DownloadIconNameKey:           "download-light-icon.svg",
		DownloadFailedIconNameKey:     "download-warning-icon.svg",
		DownloadPendingIconNameKey:    "download-wait-icon.svg",
		DownloadDeleteIconNameKey:     "download-delete-icon.svg",
		MediaDefaultIconNameKey:       "media-default-icon.svg",
		DownloadRepeatIconNameKey:     "download-repeat-icon.svg",
		DownloadSourceLinkIconNameKey: "external-link-icon.svg",
	}

	userAvatarKeysByType = map[dtypes.UserType]string{
		dtypes.UserTypeAnonymous: UserAvatarAnonymIconNameKey,
		dtypes.UserTypeGuest:     UserAvatarGuestIconNameKey,
		dtypes.UserTypeUser:      UserAvatarUserIconNameKey,
		dtypes.UserTypeAdmin:     UserAvatarAdminIconNameKey,
	}
)
