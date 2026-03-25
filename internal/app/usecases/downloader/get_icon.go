package downloader

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

func (uc *YouTubeDownloader) GetIcon(
	ctx context.Context,
	userCtx dauth.UserContext,
	fileID uuid.UUID,
) (*dmedia.ImageData, error) {
	resp, err := uc.GetFileInfo(ctx, userCtx, fileID)
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

	logo, err := uc.siteIcon.GetBySiteURL(ctx, httpx.BaseURL(resp.MediaUrl))
	if err != nil {
		return nil, err
	}

	return logo.ImageData(), nil
}
