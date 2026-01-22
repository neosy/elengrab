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
	PathGrab          = "/grab"
	PathHistory       = "/history"
	PathDownload      = "/download"
	PathFileRow       = "/file/{fileId}/row"
	PathChannelAvatar = "/channel/{channelID}/avatar"
)

func BuildPathFileRow(fileId uuid.UUID) string {
	return GroupDownloader + strings.Replace(PathFileRow, "{fileId}", fileId.String(), 1)
}

func BuildPathFileDownload(fileId uuid.UUID) string {
	return fmt.Sprintf("%s?file=%s", GroupDownloader+PathDownload, fileId)
}
