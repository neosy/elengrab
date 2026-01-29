package ytchannel

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

func (uc *YoutubeChannel) Update(ctx context.Context, channel *dmedia.YoutubeChannel) error {
	err := uc.channelRep.Update(ctx, channel)
	if err != nil {
		uc.logger.Warn("Update record error", "error", err)
		return err
	}

	if err := uc.channelCacheRep.Update(ctx, channel); err != nil {
		uc.logger.Warn("Update youtubeChannel cache error", "error", err)
		return err
	}

	return err
}
