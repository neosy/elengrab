package ytchannel

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *YoutubeChannel) FindByChannelId(ctx context.Context, channelID string) (*ddownload.YoutubeChannel, error) {
	return uc.channelRep.FindByChannelId(ctx, channelID)
}

func (uc *YoutubeChannel) ExistsByChannelId(ctx context.Context, channelID string) (bool, error) {
	return uc.channelRep.ExistsByChannelId(ctx, channelID)
}
