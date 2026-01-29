package downloader

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

func (uc *YouTubeDownloader) FindYoutubeChannelInfo(ctx context.Context, channelID string) (*dmedia.YoutubeChannel, error) {
	channel, err := uc.ytChannel.FindByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	if channel == nil {
		return nil, nil
	}

	return channel, nil
}

func (uc *YouTubeDownloader) GetYoutubeChannelInfo(ctx context.Context, channelID string) (*dmedia.YoutubeChannel, error) {
	channel, err := uc.ytChannel.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	return channel, nil
}
