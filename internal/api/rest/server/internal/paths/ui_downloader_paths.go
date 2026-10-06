package httppaths

import (
	"fmt"
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
	DownloadFilePath = "/download"
	SearchPath       = "/search"
	EventsPath       = "/events"
	MediaItemsPath   = "/items"

	// Paths Channels
	ChannelsPath = "/channels"

	// Paths items Downloader
	MediaItemPath                = MediaItemsPath + "/{itemId}"
	MediaItemRowPath             = MediaItemsPath + "/{itemId}/row"
	MediaItemDownloadRepeatPath  = MediaItemsPath + "/{itemId}/repeat"
	MediaItemImagePath           = MediaItemsPath + "/{itemId}/image"
	MediaItemMenuPath            = MediaItemsPath + "/{itemId}/menu"
	MediaItemShortLinkPath       = MediaItemsPath + "/{itemId}/short-link"
	MediaItemStreamPath          = MediaItemsPath + "/{itemId}/stream"
	MediaItemEditPath            = MediaItemsPath + "/{itemId}/edit"
	MediaItemRefreshPath         = MediaItemsPath + "/{itemId}/refresh"
	MediaItemReWatchTrackingPath = MediaItemsPath + "/{itemId}/watch-tracking"
	MediaItemWatchPositionPath   = MediaItemsPath + "/{itemId}/watch-position"

	// Paths channels Downloader
	ChannelPath = ChannelsPath + "/{channelId}"

	// Downloader Paths
	DownloaderAccountMenuPath  = DownloaderGroup + AccountMenuPath
	DownloaderSettingsMenuPath = DownloaderGroup + SettingsMenuPath
	DownloaderGrabPath         = DownloaderGroup + GrabPath
	DownloaderShareTargetPath  = DownloaderGroup + ShareTargetPath
	DownloaderItemsPath        = DownloaderGroup + MediaItemsPath
	DownloaderDownloadFilePath = DownloaderGroup + DownloadFilePath
	DownloaderSearchPath       = DownloaderGroup + SearchPath
	DownloaderEventsPath       = DownloaderGroup + EventsPath
	DownloaderChannelsPath     = DownloaderGroup + ChannelsPath
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

func BuildMediaItemEditPath(downloadID uuid.UUID) string {
	return buildMediaItemPath(MediaItemEditPath, downloadID)
}

func BuildMediaItemDownloadPath(downloadID uuid.UUID) string {
	id := idcodec.EncodeUUIDBase64URL(downloadID)
	return fmt.Sprintf("%s?itemId=%s", DownloaderGroup+DownloadFilePath, id)
}

func BuildMediaItemImagePath(downloadID uuid.UUID, verHash string, sources []dtypes.ImageSource) string {
	urlSufix := buildImageSufix(verHash, sources)
	return buildMediaItemPath(MediaItemImagePath, downloadID) + urlSufix
}

func buildChannelPath(path string, channelID uuid.UUID) string {
	id := idcodec.EncodeUUIDBase64URL(channelID)
	return DownloaderGroup + strings.Replace(path, "{channelId}", id, 1)
}

func BuildChannelPath(channelID uuid.UUID) string {
	return buildChannelPath(ChannelPath, channelID)
}
