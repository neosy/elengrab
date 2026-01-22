package downloader

import (
	"context"

	dyoutube "github.com/neosy/elengrab/internal/domain/youtube_info"
)

func (uc *YouTubeDownloader) FindYoutubeChannelInfo(ctx context.Context, channelID string) (*dyoutube.YoutubeChannel, error) {
	channel, err := uc.ytChannel.FindByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	if channel == nil {
		return nil, nil
	}

	return channel, nil
}

func (uc *YouTubeDownloader) GetYoutubeChannelInfo(ctx context.Context, channelID string) (*dyoutube.YoutubeChannel, error) {
	channel, err := uc.ytChannel.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	return channel, nil
}
