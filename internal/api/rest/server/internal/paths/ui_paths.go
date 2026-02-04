package httppaths

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// UI
const (
	// Groups
	GroupDownloader = "/ui/downloader"

	// Paths
	PathGrab               = "/grab"
	PathHistory            = "/history"
	PathDownload           = "/download"
	PathFile               = "/file/{fileId}"
	PathFileRow            = "/file/{fileId}/row"
	PathFileDownloadRepeat = "/file/{fileId}/repeat"
	PathFileLogo           = "/file/{fileId}/logo"
	PathChannelAvatar      = "/channel/{channelID}/avatar"
)

func BuildPathFile(fileID uuid.UUID) string {
	return GroupDownloader + strings.Replace(PathFile, "{fileId}", fileID.String(), 1)
}

func BuildPathFileRow(fileID uuid.UUID) string {
	return GroupDownloader + strings.Replace(PathFileRow, "{fileId}", fileID.String(), 1)
}

func BuildPathFileDownload(fileID uuid.UUID) string {
	return fmt.Sprintf("%s?file=%s", GroupDownloader+PathDownload, fileID)
}

func BuildPathFileRepeat(fileID uuid.UUID) string {
	return GroupDownloader + strings.Replace(PathFileDownloadRepeat, "{fileId}", fileID.String(), 1)
}
