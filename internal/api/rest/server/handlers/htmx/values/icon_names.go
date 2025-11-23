package htmxvalues

import (
	"os"
	"path/filepath"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const (
	DownloadIconNameKey         = "DownloadIconName"
	DownloadFailedIconNameKey   = "DownloadFailedIconName"
	DownloadPendingIconNameKey  = "DownloadPendingIconName"
	GrabResultStatusIconNameKey = "GrabResultStatusIconName"
)

var IconNames = map[string]any{
	DownloadIconNameKey:        "download-light-icon.svg",
	DownloadFailedIconNameKey:  "download-warning-icon.svg",
	DownloadPendingIconNameKey: "download-wait-icon.svg",
}

func GrabResultStatusIconName(status dtypes.FileStatus) string {
	var iconName string

	switch status {
	case dtypes.FileStatusNew, dtypes.FileStatusPending:
		iconName = IconNames[DownloadPendingIconNameKey].(string)
	case dtypes.FileStatusDone:
		iconName = IconNames[DownloadIconNameKey].(string)
	case dtypes.FileStatusFailed:
		iconName = IconNames[DownloadFailedIconNameKey].(string)
	}

	return iconName
}

func GrabResultStatusIconSvgRaw(status dtypes.FileStatus, svgDir string) string {
	filePath := filepath.Join(svgDir, GrabResultStatusIconName(status))

	data, err := os.ReadFile(filePath)
	if err != nil {
		return `<svg width="1em" height="1em"></svg>`
	}

	return string(data)
}
