package downloader

import (
	"context"

	"github.com/google/uuid"
	uchannel "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/channel"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

type Channel interface {
	Patch(
		ctx context.Context,
		channelID uuid.UUID,
		mutate func(*dmedia.Channel) error,
	) error
	UpdateChannelID(ctx context.Context, oldChannelID, newChannelID uuid.UUID) error

	IterateAll(ctx context.Context, fn func(*dmedia.Channel) error) error

	Tx(ctx context.Context, fn func(ctx context.Context) error) error
}

type channel struct {
	*uchannel.Channel
}

func (uc *downloader) Channel() Channel {
	return &channel{
		Channel: uc.channel,
	}
}
