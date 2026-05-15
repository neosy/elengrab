package icons

import (
	"html/template"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

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

	return FileRaw(iconName, svgDir)
}

func UserAvatarKeyByType(userType dtypes.UserType) string {
	key := userAvatarKeysByType[userType]
	if key == "" {
		return UserAvatarAnonymIconNameKey
	}
	return key
}
