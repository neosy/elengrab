package persistence

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type YoutubeChannelRepository interface {
	Insert(ctx context.Context, channel *ddownload.YoutubeChannel) error
	Update(ctx context.Context, channel *ddownload.YoutubeChannel) error
	FindByChannelId(ctx context.Context, channelID string) (*ddownload.YoutubeChannel, error)
	ExistsByChannelId(ctx context.Context, channelID string) (bool, error)
}
