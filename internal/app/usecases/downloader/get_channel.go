package downloader

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

func (uc *downloader) FindChannelInfo(
	ctx context.Context,
	channelID string,
	platform string,
) (*dmedia.Channel, error) {
	channel, err := uc.channel.FindByExternalChannelID(ctx, channelID, platform)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (uc *downloader) GetChannelInfo(
	ctx context.Context,
	channelID string,
	platform string,
) (*dmedia.Channel, error) {
	channel, err := uc.channel.GetByExternalChannelID(ctx, channelID, platform)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (uc *downloader) GetChannelByID(ctx context.Context, channelID uuid.UUID) (*dmedia.Channel, error) {
	return uc.channel.GetByChannelID(ctx, channelID)
}
