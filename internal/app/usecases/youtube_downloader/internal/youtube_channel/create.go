package ytchannel

import (
	"context"
	"errors"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *YoutubeChannel) Create(ctx context.Context, channel *ddownload.YoutubeChannel) error {
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
