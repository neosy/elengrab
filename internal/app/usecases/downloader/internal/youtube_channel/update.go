package ytchannel

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// Update updates an existing YouTube channel in the database.
func (uc *YoutubeChannel) Update(ctx context.Context, channel *dmedia.YoutubeChannel) error {
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
