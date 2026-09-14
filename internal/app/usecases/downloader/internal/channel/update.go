package channel

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// Update updates an existing YouTube channel in the database.
func (uc *Channel) Update(ctx context.Context, channel *dmedia.Channel) error {
	if err := channel.Validate(); err != nil {
		return err
	}

	err := uc.channelRepo().Update(ctx, channel)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return err
	}

	if err := uc.channelCacheRep.Save(ctx, channel); err != nil {
		uc.logger.Warn("Update youtubeChannel cache error", "channelID", channel.ChannelID, "error", err)
		return err
	}

	return err
}

func (uc *Channel) UpdateChannelID(ctx context.Context, oldChannelID, newChannelID uuid.UUID) error {
	return uc.channelRepo().UpdateChannelID(ctx, oldChannelID, newChannelID)
}
