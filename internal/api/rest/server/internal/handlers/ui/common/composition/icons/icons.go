package icons

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

var (
	UserMenuSearchIcon   = newIcon("UserMenuSearchIcon", "search-icon.svg")
	UserMenuDownloadIcon = newIcon("UserMenuDownloadIcon", "add-media-icon.svg")
	SearchBackArrowIcon  = newIcon("SearchBackArrowIcon", "back-arrow-icon.svg")

	UserAvatarAdminIcon  = newIcon("UserAvatarAdminIconName", "user-admin-2.svg")
	UserAvatarUserIcon   = newIcon("UserAvatarUserIconName", "user-default.svg")
	UserAvatarGuestIcon  = newIcon("UserAvatarGuestIconName", "user-guest.svg")
	UserAvatarAnonymIcon = newIcon("UserAvatarAnonymIconName", "user-anonymous-2.svg")

	IndexGrabSettingsButtonIcon = newIcon("IndexGrabSettingsButtonIconName", "settings-icon.svg")
	IndexGrabGetButtonIcon      = newIcon("IndexGrabGetButtonIconName", "download-cloud-icon.svg")

	DownloadIcon           = newIcon("DownloadIconName", "download-light-icon.svg")
	DownloadFailedIcon     = newIcon("DownloadFailedIconName", "download-warning-icon.svg")
	DownloadPendingIcon    = newIcon("DownloadPendingIconName", "download-wait-icon.svg")
	DownloadDeleteIcon     = newIcon("DownloadDeleteIconName", "download-delete-icon.svg")
	DownloadRefreshingIcon = newIcon("DownloadRefreshingIconName", "refresh-icon.svg")

	MediaPublicIcon  = newIcon("DownloadPublicIconName", "public-people.svg")
	MediaPrivateIcon = newIcon("DownloadPrivateIconName", "private-lock.svg")

	MediaAudioIcon     = newIcon("DownloadMusicIconName", "music-note-icon.svg")
	MediaShareLinkIcon = newIcon("DownloadShareLinkIconName", "link-icon.svg")
	MediaWatchedIcon   = newIcon("DownloadWatchedIconName", "check-icon.svg")

	MediaDefaultIcon                  = newIcon("MediaDefaultIconName", "media-default-icon.svg")
	DownloadRepeatIcon                = newIcon("DownloadRepeatIconName", "download-repeat-icon.svg")
	DownloadSourceLinkIcon            = newIcon("DownloadSourceLinkIconName", "external-link-icon.svg")
	DownloadMetaUserNameSeparatorIcon = newIcon("DownloadMetaUserNameSeparatorIconName", "play-icon.svg")

	CopyURLIcon = newIcon("CopyURLIconName", "copy-icon.svg")

	AccountMenuLogoutIcon = newIcon("AccountMenuLogoutIconName", "menu-logout-icon.svg")

	DownloaderRowMenuPlayIcon            = newIcon("DownloaderRowMenuPlayIconName", "menu-play-icon.svg")
	DownloaderRowMenuExternalLinkIcon    = newIcon("DownloaderRowMenuExternalLinkIconName", "menu-external-link-icon.svg")
	DownloaderRowMenuShareLinkIcon       = newIcon("DownloaderRowMenuShareLinkIconName", "menu-share-link-icon.svg")
	DownloaderRowMenuCopyLinkIcon        = newIcon("DownloaderRowMenuCopyLinkIconName", "menu-copy-link-icon.svg")
	DownloaderRowMenuDeleteLinkIcon      = newIcon("DownloaderRowMenuDeleteLinkIconName", "remove-link-icon.svg")
	DownloaderRowMenuDeleteIcon          = newIcon("DownloaderRowMenuDeleteIconName", "download-delete-icon.svg")
	DownloaderRowMenuEditIcon            = newIcon("DownloaderRowMenuEditIconName", "edit-icon.svg")
	DownloaderRowMenuUpdateMetadataIcon  = newIcon("DownloaderRowMenuUpdateMetadataIconName", "update-icon.svg")
	DownloaderRowMenuUpdateErrorInfoIcon = newIcon("DownloaderRowMenuUpdateErrorInfoIconName", "download-warning-icon.svg")

	AdminUsersIcon    = newIcon("AdminUsersIconName", "users.svg")
	AdminUserEditIcon = newIcon("AdminUserEditIconName", "user-edit.svg")
)

var (
	userAvatarIconsByType = map[dtypes.UserType]Icon{
		dtypes.UserTypeAnonymous: UserAvatarAnonymIcon,
		dtypes.UserTypeGuest:     UserAvatarGuestIcon,
		dtypes.UserTypeUser:      UserAvatarUserIcon,
		dtypes.UserTypeAdmin:     UserAvatarAdminIcon,
	}

	downloaderIconsByMediaDownloadStatus = map[dtypes.MediaDownloadStatus]Icon{
		dtypes.MediaDownloadStatusNew:    DownloadPendingIcon,
		dtypes.MediaDownloadStatusDone:   DownloadIcon,
		dtypes.MediaDownloadStatusFailed: DownloadFailedIcon,
	}
)

func UserAvatarIconByType(userType dtypes.UserType) Icon {
	return userAvatarIconsByType[userType]
}

func DownloaderIconByStatus(status dtypes.MediaDownloadStatus) Icon {
	return downloaderIconsByMediaDownloadStatus[status]
}
