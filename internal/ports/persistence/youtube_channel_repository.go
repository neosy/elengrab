package persistence

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	nmemory "github.com/neosy/elengrab/internal/pkg/ncache/memory"
)

type YoutubeChannelRepository interface {
	Insert(ctx context.Context, channel *dmedia.YoutubeChannel) error
	Update(ctx context.Context, channel *dmedia.YoutubeChannel) error
	Save(ctx context.Context, channel *dmedia.YoutubeChannel) error
	FindByChannelID(ctx context.Context, channelID string) (*dmedia.YoutubeChannel, error)
	ExistsByChannelID(ctx context.Context, channelID string) (bool, error)
}

type YoutubeChannelCacheRepository interface {
	Save(ctx context.Context, channel *dmedia.YoutubeChannel) error
	SaveNegative(ctx context.Context, channelID string) error
	FindByChannelID(ctx context.Context, channelID string) (*dmedia.YoutubeChannel, nmemory.CacheStatus, error)
	ExistsByChannelID(ctx context.Context, channelID string) (bool, error)
	CleanExpired(ctx context.Context) error
}
