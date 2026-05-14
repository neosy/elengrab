package uivalues

import (
	"html/template"
	"os"
	"path/filepath"
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
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

type iconEntry struct {
	raw template.HTML
}

func (icon *iconEntry) Copy() *iconEntry {
	return uptr.Copy(icon)
}

const iconCacheTTL = 0 * time.Hour

var (
	iconCache = memsimple.NewCacheWithDeaultCopier[string, iconEntry, *iconEntry]()

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

func IconFileNames() map[string]any {
	return iconFileNames
}

func IconFileName(key string) string {
	return iconFileNames[key].(string)
}

func IconFileRaw(fileName string, svgDir string) template.HTML {
	const svgEmpty = `<svg width="1em" height="1em"></svg>`

	icon := iconCache.Find(fileName)
	if icon != nil {
		return icon.raw
	}

	filePath := filepath.Join(svgDir, fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return svgEmpty
	}

	iconCache.Save(fileName, &iconEntry{template.HTML(data)}, iconCacheTTL)

	return template.HTML(data)
}

func IconFileRawByKey(key string, svgDir string) template.HTML {
	return IconFileRaw(IconFileName(key), svgDir)
}

func DownloaderResultStatusIconSvgRaw(status dtypes.MediaDownloadStatus, svgDir string) template.HTML {
	var iconName string

	switch status {
	case dtypes.MediaDownloadStatusNew, dtypes.MediaDownloadStatusPending:
		iconName = iconFileNames[DownloadPendingIconNameKey].(string)
	case dtypes.MediaDownloadStatusDone:
		iconName = iconFileNames[DownloadIconNameKey].(string)
	case dtypes.MediaDownloadStatusFailed:
		iconName = iconFileNames[DownloadFailedIconNameKey].(string)
	}

	return IconFileRaw(iconName, svgDir)
}

func UserAvatarKeyByType(userType dtypes.UserType) string {
	key := userAvatarKeysByType[userType]
	if key == "" {
		return UserAvatarAnonymIconNameKey
	}
	return key
}
