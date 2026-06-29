package httppaths

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

// UI
const (
	// Groups
	GroupDownloader = "/downloader"
	GroupAccount    = "/account"

	// Paths Downloader
	PathAccountMenu  = "/account-menu"
	PathSettingsMenu = "/settings-menu"
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
	PathMediaItemEdit           = "/items/{itemId}/edit"
	PathMediaItemRefresh        = "/items/{itemId}/refresh"

	// Paths channels in Downloader
	PathChannelAvatar = "/channels/{channelId}/avatar"

	// Paths Account
	PathRegister = "/register"
	PathLogin    = "/login"
	PathLogout   = "/logout"

	// Paths short links, e.g. /s/{shortCode}
	PathShortLink       = "/{shortCode}"
	PathStreamShortCode = "/stream/{shortCode}"
)

func buildMediaItemPath(path string, downloadID uuid.UUID) string {
	id := idcodec.EncodeUUIDBase64URL(downloadID)
	return GroupDownloader + strings.Replace(path, "{itemId}", id, 1)
}

func BuildMediaItemPath(downloadID uuid.UUID) string {
	return buildMediaItemPath(PathMediaItem, downloadID)
}

func BuildPathMediaItemRow(downloadID uuid.UUID) string {
	return buildMediaItemPath(PathMediaItemRow, downloadID)
}

func BuildPathMediaItemDownloadRepeat(downloadID uuid.UUID) string {
	return buildMediaItemPath(PathMediaItemDownloadRepeat, downloadID)
}

func BuildPathMediaItemStream(downloadID uuid.UUID) string {
	return buildMediaItemPath(PathMediaItemStream, downloadID)
}

func BuildPathMediaItemWatch(downloadID uuid.UUID) string {
	return buildMediaItemPath(PathMediaItemWatch, downloadID)
}

func BuildPathMediaItemEdit(downloadID uuid.UUID) string {
	return buildMediaItemPath(PathMediaItemEdit, downloadID)
}

func BuildPathStreamShortCode(shortCode string) string {
	return GroupDownloader + strings.Replace(PathStreamShortCode, "{shortCode}", shortCode, 1)
}

func BuildPathMediaItemDownload(downloadID uuid.UUID) string {
	id := idcodec.EncodeUUIDBase64URL(downloadID)
	return fmt.Sprintf("%s?itemId=%s", GroupDownloader+PathDownloadFile, id)
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

	return buildMediaItemPath(PathMediaItemImage, downloadID) + urlSufix
}
