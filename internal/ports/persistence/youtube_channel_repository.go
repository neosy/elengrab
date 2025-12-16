package persistence

import (
	"context"

	dyoutube "github.com/neosy/elengrab/internal/domain/youtube_info"
)

type YoutubeChannelRepository interface {
	Insert(ctx context.Context, channel *dyoutube.YoutubeChannel) error
	Update(ctx context.Context, channel *dyoutube.YoutubeChannel) error
	Save(ctx context.Context, channel *dyoutube.YoutubeChannel) error
	FindByChannelID(ctx context.Context, channelID string) (*dyoutube.YoutubeChannel, error)
	ExistsByChannelID(ctx context.Context, channelID string) (bool, error)
}

type YoutubeChannelCacheRepository interface {
	YoutubeChannelRepository
	CleanExpired(ctx context.Context) error
}
