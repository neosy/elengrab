package downloader

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/pkg/httpx"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (uc *YouTubeDownloader) GetLogo(ctx context.Context, userID uuid.UUID, fileID uuid.UUID) (*dmedia.ImageData, error) {
	resp, err := uc.GetFileInfo(ctx, userID, fileID)
	if err != nil {
		return nil, err
	}

	if resp.YoutubeChannelID != nil {
		channel, _ := uc.FindYoutubeChannelInfo(ctx, *resp.YoutubeChannelID)
		if channel != nil && len(channel.ImageRaw) > 0 {
			return &dmedia.ImageData{
				Raw:    channel.ImageRaw,
				Format: channel.ImageFormat,
			}, nil
		}
	}

	logo, err := uc.siteLogo.GetBySiteURL(ctx, httpx.GetBaseURL(resp.YoutubeUrl))
	if err != nil {
		return nil, err
	}

	return uptr.Any(logo.ImageData()), nil
}
