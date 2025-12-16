package ytchannel

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type YoutubeChannel struct {
	logger *slog.Logger

	// repositories
	channelRep persistence.YoutubeChannelRepository

	// cache
	channelCacheRep persistence.YoutubeChannelRepository
}

func NewYoutubeChannel(
	logger *slog.Logger,
	channelRep persistence.YoutubeChannelRepository,
	channelCacheRep persistence.YoutubeChannelRepository,
) *YoutubeChannel {
	return &YoutubeChannel{
		logger:          logger,
		channelRep:      channelRep,
		channelCacheRep: channelCacheRep,
	}
}
