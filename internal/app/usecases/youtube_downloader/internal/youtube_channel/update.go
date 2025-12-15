package ytchannel

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *YoutubeChannel) Update(ctx context.Context, channel *ddownload.YoutubeChannel) error {
	err := uc.channelRep.Update(ctx, channel)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return err
	}
	return err
}
