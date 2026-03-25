package uivalues

import (
	"os"
	"path/filepath"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const (
	UserAvatarAdminIconNameKey  = "UserAvatarAdminIconName"
	UserAvatarUserIconNameKey   = "UserAvatarUserIconName"
	UserAvatarGuestIconNameKey  = "UserAvatarGuestIconName"
	UserAvatarAnonymIconNameKey = "UserAvatarAnonymIconName"
	DownloadIconNameKey         = "DownloadIconName"
	DownloadFailedIconNameKey   = "DownloadFailedIconName"
	DownloadPendingIconNameKey  = "DownloadPendingIconName"
	GrabResultStatusIconNameKey = "GrabResultStatusIconName"
	DownloadDeleteIconNameKey   = "DownloadDeleteIconName"
	MediaDefaultIconNameKey     = "MediaDefaultIconName"
	DownloadRepeatIconNameKey   = "DownloadRepeatIconName"
)

var iconFileNames = map[string]any{
	UserAvatarAdminIconNameKey:  "user-admin-2.svg",
	UserAvatarUserIconNameKey:   "user-guest.svg",
	UserAvatarGuestIconNameKey:  "user-guest.svg",
	UserAvatarAnonymIconNameKey: "user-anonymous-2.svg",
	DownloadIconNameKey:         "download-light-icon.svg",
	DownloadFailedIconNameKey:   "download-warning-icon.svg",
	DownloadPendingIconNameKey:  "download-wait-icon.svg",
	DownloadDeleteIconNameKey:   "download-delete-icon.svg",
	MediaDefaultIconNameKey:     "media-default-icon.svg",
	DownloadRepeatIconNameKey:   "download-repeat-icon.svg",
}

func IconFileNames() map[string]any {
	return iconFileNames
}

func IconFileName(key string) string {
	return iconFileNames[key].(string)
}

func DownloadResultStatusIconFileName(status dtypes.FileStatus) string {
	var iconName string

	switch status {
	case dtypes.FileStatusNew, dtypes.FileStatusPending:
		iconName = iconFileNames[DownloadPendingIconNameKey].(string)
	case dtypes.FileStatusDone:
		iconName = iconFileNames[DownloadIconNameKey].(string)
	case dtypes.FileStatusFailed:
		iconName = iconFileNames[DownloadFailedIconNameKey].(string)
	}

	return iconName
}

func IconFileRaw(fileName string, svgDir string) string {
	filePath := filepath.Join(svgDir, fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return `<svg width="1em" height="1em"></svg>`
	}

	return string(data)
}

func IconFileRawByKey(key string, svgDir string) string {
	return IconFileRaw(IconFileName(key), svgDir)
}

func DownloadResultStatusIconSvgRaw(status dtypes.FileStatus, svgDir string) string {
	return IconFileRaw(DownloadResultStatusIconFileName(status), svgDir)
}
