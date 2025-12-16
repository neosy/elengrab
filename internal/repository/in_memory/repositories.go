package inmemory

import (
	"time"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

// Repositories groups all database repositories.
type Repositories struct {
	DownloadState  persistence.DownloadStateCacheRepository
	YoutubeChannel persistence.YoutubeChannelCacheRepository
}

type Dependencies struct {
	DownloadStateTTL  time.Duration
	YoutubeChannelTTL time.Duration
}

// New returns a new Repositories struct.
func New(deps Dependencies) *Repositories {
	return &Repositories{
		DownloadState:  newDownloadStateRepository(deps.DownloadStateTTL),
		YoutubeChannel: newYoutubeChannelRepository(deps.YoutubeChannelTTL),
	}
}
