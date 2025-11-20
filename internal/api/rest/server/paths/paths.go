package httppaths

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Index
const (
	PathIndex = "/"
)

// UI
const (
	// Groups
	GroupStatic     = "/ui/static"
	GroupDownloader = "/ui/downloader"

	// Paths
	PathGrab     = "/grab"
	PathHistory  = "/history"
	PathDownload = "/download"
	PathFileRow  = "/file/{fileId}/row"
	PathCssFiles = "/css/{filepath:*}"
	PathImgFiles = "/img/{filepath:*}"
	PathJsFiles  = "/js/{filepath:*}"
	PathPwaFiles = "/pwa/{filepath:*}"
)

func BuildPathFileRow(fileId uuid.UUID) string {
	return GroupDownloader + strings.Replace(PathFileRow, "{fileId}", fileId.String(), 1)
}

func BuildPathFileDownload(fileId uuid.UUID) string {
	return fmt.Sprintf("%s?file=%s", GroupDownloader+PathDownload, fileId)
}
