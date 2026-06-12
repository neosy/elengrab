package inmemory

import (
	"time"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

// Repositories groups all database repositories.
type Repositories struct {
	DownloadState  persistence.DownloadStateCacheRepository
	YoutubeChannel persistence.YoutubeChannelCacheRepository
	SiteLogo       persistence.SiteLogoCacheRepository
	Thumbnail      persistence.ThumbnailCacheRepository
	ThumbnailFile  persistence.ThumbnailFileCacheRepository
	AssetFile      persistence.AssetFileCacheRepository
}

type Dependencies struct {
	DownloadStateCacheTTL  time.Duration
	YoutubeChannelCacheTTL time.Duration
	SiteLogoCacheTTL       time.Duration
	ThumbnailCacheTTL      time.Duration
	ThumbnailFileCacheTTL  time.Duration
	AssetFileCacheTTL      time.Duration
}

// New returns a new Repositories struct.
func New(deps Dependencies) *Repositories {
	return &Repositories{
		DownloadState:  newDownloadStateRepository(deps.DownloadStateCacheTTL),
		YoutubeChannel: newYoutubeChannelRepository(deps.YoutubeChannelCacheTTL),
		SiteLogo:       newSiteLogoRepository(deps.SiteLogoCacheTTL),
		Thumbnail:      newThumbnailRepository(deps.ThumbnailCacheTTL),
		ThumbnailFile:  newThumbnailFileRepository(deps.ThumbnailFileCacheTTL),
		AssetFile:      newAssetFileRepository(deps.AssetFileCacheTTL),
	}
}
