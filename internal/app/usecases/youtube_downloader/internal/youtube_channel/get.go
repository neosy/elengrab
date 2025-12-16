package ytchannel

import (
	"context"

	dyoutube "github.com/neosy/elengrab/internal/domain/youtube_info"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
)

// FindByChannelID
// Channel may not exist — caller decides what to do
func (uc *YoutubeChannel) FindByChannelID(ctx context.Context, channelID string) (*dyoutube.YoutubeChannel, error) {
	if channelID == "" {
		return nil, nil
	}

	channel, _ := uc.channelCacheRep.FindByChannelID(ctx, channelID)
	if channel != nil {
		return channel, nil
	}

	channel, err := uc.channelRep.FindByChannelID(ctx, channelID)
	if err != nil {
		uc.logger.Warn("Failed get youtubeChannel", "error", err)
		return nil, errorx.NewByErr(err, exceptionx.ERROR)
	}

	err = uc.channelCacheRep.Insert(ctx, channel)
	if err != nil {
		uc.logger.Warn("Failed to insert youtubeChannel cache", "error", err)
	}

	return channel, nil
}

// GetByChannelID
// Channel MUST exist — otherwise NOT_FOUND
func (uc *YoutubeChannel) GetByChannelID(ctx context.Context, channelID string) (*dyoutube.YoutubeChannel, error) {
	channel, err := uc.FindByChannelID(ctx, channelID)
	if err != nil {
		return nil, errorx.NewByErr(err, exceptionx.ERROR)
	}

	if channel == nil {
		uc.logger.Warn("Channel not found", "channelId", channelID)
		return nil, errorx.New("channel not found", exceptionx.NOT_FOUND)
	}

	return channel, nil
}

func (uc *YoutubeChannel) ExistsByChannelID(ctx context.Context, channelID string) (bool, error) {
	exists, _ := uc.channelCacheRep.ExistsByChannelID(ctx, channelID)
	if exists {
		return exists, nil
	}

	exists, err := uc.channelRep.ExistsByChannelID(ctx, channelID)
	if err != nil {
		uc.logger.Warn("Failed to check if YouTube channel exists", "channelId", channelID, "error", err)
	}

	return exists, nil
}
