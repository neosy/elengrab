package ytchannel

import (
	"context"
	"errors"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// Create creates a new YouTube channel in the database.
func (uc *YoutubeChannel) Create(ctx context.Context, channel *dmedia.YoutubeChannel) error {
	if channel == nil {
		uc.logger.Warn("Nil pointer in function")
		return errors.New("function parameter is a null pointer")
	}

	err := uc.channelRep.Insert(ctx, channel)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return err
	}

	channel, _ = uc.channelRep.FindByChannelID(ctx, channel.ChannelID)
	if channel != nil {
		err := uc.channelCacheRep.Save(ctx, channel)
		if err != nil {
			uc.logger.Warn(
				"Failed to save cache",
				"channelID", channel.ChannelID,
				"error", err,
			)
		}
	}

	return nil
}
