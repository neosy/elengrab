package channel

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

func (uc *Channel) Patch(
	ctx context.Context,
	channelID uuid.UUID,
	mutate func(*dmedia.Channel) error,
) error {
	err := uc.Tx(ctx, func(ctx context.Context) error {
		channel, err := uc.FindByChannelIDNoCache(ctx, channelID)
		if err != nil {
			return err
		}

		if err := mutate(channel); err != nil {
			return err
		}

		err = uc.Update(ctx, channel)
		if err != nil {
			return err
		}

		uc.channelCacheRep.Delete(ctx, channel.ChannelID)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
