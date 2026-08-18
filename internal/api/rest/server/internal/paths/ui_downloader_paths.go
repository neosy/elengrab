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
	DownloaderGroup = "/downloader"

	// Paths Downloader
	AccountMenuPath  = "/account-menu"
	SettingsMenuPath = "/settings-menu"
	GrabPath         = "/grab"
	ShareTargetPath  = "/share-target"
	HistoryPath      = "/history"
	DownloadFilePath = "/download"
	SearchPath       = "/search"
	EventsPath       = "/events"

	// Downloader Paths
	DownloaderAccountMenuPath  = DownloaderGroup + AccountMenuPath
	DownloaderSettingsMenuPath = DownloaderGroup + SettingsMenuPath
	DownloaderGrabPath         = DownloaderGroup + GrabPath
	DownloaderShareTargetPath  = DownloaderGroup + ShareTargetPath
	DownloaderHistoryPath      = DownloaderGroup + HistoryPath
	DownloaderDownloadFilePath = DownloaderGroup + DownloadFilePath
	DownloaderSearchPath       = DownloaderGroup + SearchPath
	DownloaderEventsPath       = DownloaderGroup + EventsPath

	// Paths items Downloader
	MediaItemPath                = "/items/{itemId}"
	MediaItemRowPath             = "/items/{itemId}/row"
	MediaItemDownloadRepeatPath  = "/items/{itemId}/repeat"
	MediaItemImagePath           = "/items/{itemId}/image"
	MediaItemMenuPath            = "/items/{itemId}/menu"
	MediaItemShortLinkPath       = "/items/{itemId}/short-link"
	MediaItemStreamPath          = "/items/{itemId}/stream"
	MediaItemWatchPath           = "/items/{itemId}/watch"
	MediaItemEditPath            = "/items/{itemId}/edit"
	MediaItemRefreshPath         = "/items/{itemId}/refresh"
	MediaItemReWatchTrackingPath = "/items/{itemId}/watch-tracking"
	MediaItemWatchPositionPath   = "/items/{itemId}/watch-position"

	// Paths channels in Downloader
	ChannelAvatarPath = "/channels/{channelId}/avatar"

	// Paths short links, e.g. /s/{shortCode}
	ShortLinkPath       = "/{shortCode}"
	StreamShortCodePath = "/stream/{shortCode}"
)

func buildMediaItemPath(path string, downloadID uuid.UUID) string {
	id := idcodec.EncodeUUIDBase64URL(downloadID)
	return DownloaderGroup + strings.Replace(path, "{itemId}", id, 1)
}

func BuildMediaItemPath(downloadID uuid.UUID) string {
	return buildMediaItemPath(MediaItemPath, downloadID)
}

func BuildMediaItemRowPath(downloadID uuid.UUID) string {
	return buildMediaItemPath(MediaItemRowPath, downloadID)
}

func BuildMediaItemDownloadRepeatPath(downloadID uuid.UUID) string {
	return buildMediaItemPath(MediaItemDownloadRepeatPath, downloadID)
}

func BuildMediaItemStreamPath(downloadID uuid.UUID) string {
	return buildMediaItemPath(MediaItemStreamPath, downloadID)
}

func BuildMediaItemWatchPath(downloadID uuid.UUID) string {
	return buildMediaItemPath(MediaItemWatchPath, downloadID)
}

func BuildMediaItemEditPath(downloadID uuid.UUID) string {
	return buildMediaItemPath(MediaItemEditPath, downloadID)
}

func BuildStreamShortCodePath(shortCode string) string {
	return DownloaderGroup + strings.Replace(StreamShortCodePath, "{shortCode}", shortCode, 1)
}

func BuildMediaItemDownloadPath(downloadID uuid.UUID) string {
	id := idcodec.EncodeUUIDBase64URL(downloadID)
	return fmt.Sprintf("%s?itemId=%s", DownloaderGroup+DownloadFilePath, id)
}

func BuildMediaItemImagePath(downloadID uuid.UUID, verHash string, sources []dtypes.ImageSource) string {
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

	return buildMediaItemPath(MediaItemImagePath, downloadID) + urlSufix
}
