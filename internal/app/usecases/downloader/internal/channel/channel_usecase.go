package channel

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type Channel struct {
	logger *slog.Logger

	// repositories
	channelRepo persistence.ChannelRepositoryFactory

	// cache
	channelCacheRep persistence.ChannelCacheRepository
}

func NewChannel(
	logger *slog.Logger,
	channelRepo persistence.ChannelRepositoryFactory,
	channelCacheRep persistence.ChannelCacheRepository,
) *Channel {
	return &Channel{
		logger:          logger,
		channelRepo:     channelRepo,
		channelCacheRep: channelCacheRep,
	}
}
