package httppaths

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// UI
const (
	// Groups
	GroupDownloader = "/downloader"
	GroupAccount    = "/account"

	// Paths Downloader
	PathGrab        = "/grab"
	PathShareTarget = "/share-target"
	PathHistory     = "/history"
	PathDownload    = "/download"
	PathSearch      = "/search"
	PathAccountMenu = "/account-menu"

	PathFilesEvents        = "/files/events"
	PathFile               = "/file/{fileId}"
	PathFileRow            = "/file/{fileId}/row"
	PathFileDownloadRepeat = "/file/{fileId}/repeat"
	PathFileLogo           = "/file/{fileId}/logo"
	PathFileMenu           = "/files/{fileId}/menu"
	PathFileShortLink      = "/files/{fileId}/short-link"
	PathFileStream         = "/files/{fileId}/stream"
	PathFileWatch          = "/file/{fileId}/watch"

	PathChannelAvatar = "/channel/{channelID}/avatar"

	// Paths Account
	PathRegister = "/register"
	PathLogin    = "/login"
	PathLogout   = "/logout"

	// Short link
	PathShortLink       = "/{shortCode}"
	PathStreamShortCode = "/stream/{shortCode}"
)

func BuildPathFile(path string, fileID uuid.UUID) string {
	return GroupDownloader + strings.Replace(path, "{fileId}", fileID.String(), 1)
}

func BuildPathFileRow(fileID uuid.UUID) string {
	return BuildPathFile(PathFileRow, fileID)
}

func BuildPathFileRepeat(fileID uuid.UUID) string {
	return BuildPathFile(PathFileDownloadRepeat, fileID)
}

func BuildPathFileStream(fileID uuid.UUID) string {
	return BuildPathFile(PathFileStream, fileID)
}

func BuildPathFileWatch(fileID uuid.UUID) string {
	return BuildPathFile(PathFileWatch, fileID)
}

func BuildPathStreamShortCode(shortCode string) string {
	return GroupDownloader + strings.Replace(PathStreamShortCode, "{shortCode}", shortCode, 1)
}

func BuildPathFileDownload(fileID uuid.UUID) string {
	return fmt.Sprintf("%s?file=%s", GroupDownloader+PathDownload, fileID)
}
