package channel

import (
	"context"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// Create creates a new YouTube channel in the database.
func (uc *Channel) Create(ctx context.Context, channel *dmedia.Channel) error {
	if channel == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	if channel.ChannelID == uuid.Nil {
		channel.ChannelID = uuid.New()
	}

	if err := channel.Validate(); err != nil {
		return err
	}

	defer uc.channelCacheRep.Delete(ctx, channel.ChannelID)

	err := uc.channelRepo().Insert(ctx, channel)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return err
	}

	return nil
}
