package channel

import (
	"context"

	"github.com/google/uuid"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *Channel) FindByChannelIDNoCache(ctx context.Context, channelID uuid.UUID) (*dmedia.Channel, error) {
	channel, err := uc.channelRepo().FindByChannelID(ctx, channelID)
	if err != nil {
		uc.logger.Warn("Failed get youtubeChannel", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return channel, nil
}

// FindByChannelID
// Channel may not exist — caller decides what to do
func (uc *Channel) FindByChannelID(
	ctx context.Context,
	channelID uuid.UUID,
) (*dmedia.Channel, error) {
	if channelID == uuid.Nil {
		return nil, nil
	}

	channel, cacheStatus, _ := uc.channelCacheRep.FindByChannelID(ctx, channelID)
	if channel != nil {
		return channel, nil
	}
	if cacheStatus == memsimple.CacheStatusNegativeHit {
		return nil, nil
	}

	channel, err := uc.FindByChannelIDNoCache(ctx, channelID)
	if err != nil {
		return nil, err
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
func (uc *Channel) GetByChannelID(ctx context.Context, channelID uuid.UUID) (*dmedia.Channel, error) {
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

func (uc *Channel) ExistsByChannelID(ctx context.Context, channelID uuid.UUID) (bool, error) {
	exists, _ := uc.channelCacheRep.ExistsByChannelID(ctx, channelID)
	if exists {
		return exists, nil
	}

	exists, err := uc.channelRepo().ExistsByChannelID(ctx, channelID)
	if err != nil {
		uc.logger.Warn("Failed to check if YouTube channel exists", "channelId", channelID, "error", err)
	}

	return exists, nil
}

func (uc *Channel) FindByExternalChannelIDNoCache(
	ctx context.Context,
	externalID string,
	platform string,
) (*dmedia.Channel, error) {
	if externalID == "" {
		return nil, nil
	}

	channel, err := uc.channelRepo().FindByExternalChannelID(ctx, externalID, platform)
	if err != nil {
		uc.logger.Warn("Failed get youtubeChannel", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return channel, nil
}

func (uc *Channel) FindByExternalChannelID(
	ctx context.Context,
	externalID string,
	platform string,
) (*dmedia.Channel, error) {
	if externalID == "" {
		return nil, nil
	}

	channel, cacheStatus, _ := uc.channelCacheRep.FindByExternalChannelID(ctx, externalID, platform)
	if channel != nil {
		return channel, nil
	}
	if cacheStatus == memsimple.CacheStatusNegativeHit {
		return nil, nil
	}

	channel, err := uc.FindByExternalChannelIDNoCache(ctx, externalID, platform)
	if err != nil {
		return nil, err
	}

	if channel == nil {
		err := uc.channelCacheRep.SaveNegativeByExternalChannelID(ctx, externalID, platform)
		if err != nil {
			uc.logger.Warn(
				"Failed to insert channel cache",
				"externalChannelID", externalID,
				"platform", platform,
				"error", err,
			)
		}
		return nil, nil
	}

	uc.channelCacheRep.Save(ctx, channel)

	return channel, nil
}

func (uc *Channel) GetByExternalChannelID(
	ctx context.Context,
	channelID string,
	platform string,
) (*dmedia.Channel, error) {
	channel, err := uc.FindByExternalChannelID(ctx, channelID, platform)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if channel == nil {
		uc.logger.Warn(
			"Channel not found",
			"externalChannelId", channelID,
			"platform", platform,
		)
		return nil, errorx.New("channel not found", exceptionx.NOT_FOUND)
	}

	return channel, nil
}

func (uc *Channel) GetAllIDs(ctx context.Context) ([]uuid.UUID, error) {
	return uc.channelRepo().GetIDs(ctx)
}

func (uc *Channel) IterateAll(ctx context.Context, fn func(*dmedia.Channel) error) error {
	return uc.channelRepo().IterateAll(ctx, fn)
}
