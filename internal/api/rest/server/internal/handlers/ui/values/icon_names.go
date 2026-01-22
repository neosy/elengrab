package uivalues

import (
	"os"
	"path/filepath"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const (
	DownloadIconNameKey              = "DownloadIconName"
	DownloadFailedIconNameKey        = "DownloadFailedIconName"
	DownloadPendingIconNameKey       = "DownloadPendingIconName"
	GrabResultStatusIconNameKey      = "GrabResultStatusIconName"
	DownloadDeleteIconNameKey        = "DownloadDeleteIconName"
	YoutubeChannelDefaultIconNameKey = "YoutubeChannelDefaultIconName"
)

var iconFileNames = map[string]any{
	DownloadIconNameKey:              "download-light-icon.svg",
	DownloadFailedIconNameKey:        "download-warning-icon.svg",
	DownloadPendingIconNameKey:       "download-wait-icon.svg",
	DownloadDeleteIconNameKey:        "download-delete-icon.svg",
	YoutubeChannelDefaultIconNameKey: "youtube_channel_default_avatar-icon_3.svg",
}

func IconFileNames() map[string]any {
	return iconFileNames
}

func IconFileName(key string) string {
	return iconFileNames[key].(string)
}

func FrabResultStatusIconFileName(status dtypes.FileStatus) string {
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

func GrabResultStatusIconSvgRaw(status dtypes.FileStatus, svgDir string) string {
	return IconFileRaw(FrabResultStatusIconFileName(status), svgDir)
}
