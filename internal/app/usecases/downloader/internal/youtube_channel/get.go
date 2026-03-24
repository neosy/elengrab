package ytchannel

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
	nmemory "github.com/neosy/elengrab/pkg/ncache/memory"
)

// FindByChannelID
// Channel may not exist — caller decides what to do
func (uc *YoutubeChannel) FindByChannelID(ctx context.Context, channelID string) (*dmedia.YoutubeChannel, error) {
	if channelID == "" {
		return nil, nil
	}

	channel, cacheStatus, _ := uc.channelCacheRep.FindByChannelID(ctx, channelID)
	if channel != nil {
		return channel, nil
	}
	if cacheStatus == nmemory.CacheStatusNegativeHit {
		return nil, nil
	}

	channel, err := uc.channelRep.FindByChannelID(ctx, channelID)
	if err != nil {
		uc.logger.Warn("Failed get youtubeChannel", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if channel == nil {
		err := uc.channelCacheRep.SaveNegative(ctx, channelID)
		if err != nil {
			uc.logger.Warn("Failed to insert youtubeChannel cache", "channelID", channelID, "error", err)
		}
		return nil, nil
	}

	err = uc.channelCacheRep.Save(ctx, channel)
	if err != nil {
		uc.logger.Warn("Failed to insert youtubeChannel cache", "channelID", channelID, "error", err)
	}

	return channel, nil
}

// GetByChannelID
// Channel MUST exist — otherwise NOT_FOUND
func (uc *YoutubeChannel) GetByChannelID(ctx context.Context, channelID string) (*dmedia.YoutubeChannel, error) {
	channel, err := uc.FindByChannelID(ctx, channelID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
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
