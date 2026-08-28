package ytchannel

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type YoutubeChannel struct {
	logger *slog.Logger

	// repositories
	channelRepo persistence.YoutubeChannelRepositoryFactory

	// cache
	channelCacheRep persistence.YoutubeChannelCacheRepository
}

func NewYoutubeChannel(
	logger *slog.Logger,
	channelRepo persistence.YoutubeChannelRepositoryFactory,
	channelCacheRep persistence.YoutubeChannelCacheRepository,
) *YoutubeChannel {
	return &YoutubeChannel{
		logger:          logger,
		channelRepo:     channelRepo,
		channelCacheRep: channelCacheRep,
	}
}
