package inmemory

import (
	"time"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

// Repositories groups all database repositories.
type Repositories struct {
	MediaDownload      persistence.MediaDownloadCacheRepository
	DownloadState      persistence.DownloadStateCacheRepository
	MediaWatchStat     persistence.MediaWatchStatCacheRepository
	MediaWatchPosition persistence.MediaWatchPositionCacheRepository
	YoutubeChannel     persistence.YoutubeChannelCacheRepository
	SiteLogo           persistence.SiteLogoCacheRepository
	Thumbnail          persistence.ThumbnailCacheRepository
	ThumbnailFile      persistence.ThumbnailFileCacheRepository
	AssetFile          persistence.AssetFileCacheRepository
}

type Dependencies struct {
	MediaDownloadCacheTTL      time.Duration
	DownloadStateCacheTTL      time.Duration
	MediaWatchStatCacheTTL     time.Duration
	MediaWatchPositionCacheTTL time.Duration
	YoutubeChannelCacheTTL     time.Duration
	SiteLogoCacheTTL           time.Duration
	ThumbnailCacheTTL          time.Duration
	ThumbnailFileCacheTTL      time.Duration
	AssetFileCacheTTL          time.Duration
}

// New returns a new Repositories struct.
func New(deps Dependencies) *Repositories {
	return &Repositories{
		MediaDownload:      newMediaDownloadRepository(deps.MediaDownloadCacheTTL),
		DownloadState:      newDownloadStateRepository(deps.DownloadStateCacheTTL),
		MediaWatchStat:     newMediaWatchStatRepository(deps.MediaWatchStatCacheTTL),
		MediaWatchPosition: newMediaWatchPositionRepository(deps.MediaWatchPositionCacheTTL),
		YoutubeChannel:     newYoutubeChannelRepository(deps.YoutubeChannelCacheTTL),
		SiteLogo:           newSiteLogoRepository(deps.SiteLogoCacheTTL),
		Thumbnail:          newThumbnailRepository(deps.ThumbnailCacheTTL),
		ThumbnailFile:      newThumbnailFileRepository(deps.ThumbnailFileCacheTTL),
		AssetFile:          newAssetFileRepository(deps.AssetFileCacheTTL),
	}
}
