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

	return nil
}
