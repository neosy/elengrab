package httppaths

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// UI
const (
	// Groups
	GroupDownloader = "/downloader"
	GroupAccount    = "/account"

	// Paths Downloader
	PathAccountMenu  = "/account-menu"
	PathGrab         = "/grab"
	PathShareTarget  = "/share-target"
	PathHistory      = "/history"
	PathDownloadFile = "/download"
	PathSearch       = "/search"
	PathEvents       = "/events"

	// Paths items Downloader
	PathMediaItem               = "/items/{itemId}"
	PathMediaItemRow            = "/items/{itemId}/row"
	PathMediaItemDownloadRepeat = "/items/{itemId}/repeat"
	PathMediaItemImage          = "/items/{itemId}/image"
	PathMediaItemMenu           = "/items/{itemId}/menu"
	PathMediaItemShortLink      = "/items/{itemId}/short-link"
	PathMediaItemStream         = "/items/{itemId}/stream"
	PathMediaItemWatch          = "/items/{itemId}/watch"

	PathChannelAvatar = "/channels/{channelId}/avatar"

	// Paths Account
	PathRegister = "/register"
	PathLogin    = "/login"
	PathLogout   = "/logout"

	// Short link
	PathShortLink       = "/{shortCode}"
	PathStreamShortCode = "/stream/{shortCode}"
)

func BuildPathMediaItem(path string, downloadID uuid.UUID) string {
	return GroupDownloader + strings.Replace(path, "{itemId}", downloadID.String(), 1)
}

func BuildPathMediaItemRow(downloadID uuid.UUID) string {
	return BuildPathMediaItem(PathMediaItemRow, downloadID)
}

func BuildPathMediaItemDownloadRepeat(downloadID uuid.UUID) string {
	return BuildPathMediaItem(PathMediaItemDownloadRepeat, downloadID)
}

func BuildPathMediaItemStream(downloadID uuid.UUID) string {
	return BuildPathMediaItem(PathMediaItemStream, downloadID)
}

func BuildPathMediaItemWatch(downloadID uuid.UUID) string {
	return BuildPathMediaItem(PathMediaItemWatch, downloadID)
}

func BuildPathStreamShortCode(shortCode string) string {
	return GroupDownloader + strings.Replace(PathStreamShortCode, "{shortCode}", shortCode, 1)
}

func BuildPathMediaItemDownload(downloadID uuid.UUID) string {
	return fmt.Sprintf("%s?itemId=%s", GroupDownloader+PathDownloadFile, downloadID)
}

func BuildPathMediaItemImage(downloadID uuid.UUID, verHash string, sources []dtypes.ImageSource) string {
	var sourceStrings []string
	for _, src := range sources {
		if src.Exists() {
			sourceStrings = append(sourceStrings, src.String())
		}
	}

	var urlValues url.Values
	if len(sourceStrings) > 0 {
		urlValues = url.Values{}
		urlValues.Set("source", strings.Join(sourceStrings, ","))
	}

	if verHash != "" {
		urlValues.Set("v", verHash)
	}

	var urlSufix string
	if len(urlValues) > 0 {
		urlSufix = "?" + urlValues.Encode()
	}

	return BuildPathMediaItem(PathMediaItemImage, downloadID) + urlSufix
}
