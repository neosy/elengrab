package persistence

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/internal/pkg/cache/memory"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type YoutubeChannelRepositoryFactory func() YoutubeChannelRepository

type YoutubeChannelRepository interface {
	Transactional

	Insert(ctx context.Context, channel *dmedia.YoutubeChannel) error
	Update(ctx context.Context, channel *dmedia.YoutubeChannel) error
	Save(ctx context.Context, channel *dmedia.YoutubeChannel) error
	FindByChannelID(ctx context.Context, channelID string) (*dmedia.YoutubeChannel, error)
	ExistsByChannelID(ctx context.Context, channelID string) (bool, error)
}

type YoutubeChannelCacheRepository interface {
	memory.CacheRepository

	Save(ctx context.Context, channel *dmedia.YoutubeChannel) error
	SaveNegative(ctx context.Context, channelID string) error

	FindByChannelID(ctx context.Context, channelID string) (*dmedia.YoutubeChannel, memsimple.CacheStatus, error)
	ExistsByChannelID(ctx context.Context, channelID string) (bool, error)

	CleanExpired(ctx context.Context) error
}
