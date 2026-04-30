package downloader

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

func (uc *Downloader) GetFileLogo(
	ctx context.Context,
	userCtx dauth.UserContext,
	fileID uuid.UUID,
) (*dtypes.ImageData, error) {
	resp, err := uc.GetFileInfo(ctx, userCtx, fileID)
	if err != nil {
		return nil, err
	}

	if resp.MediaInfo != nil && resp.MediaInfo.GetThumbnailID() != nil {
		thumbnail, _ := uc.thumbnail.GetByThumbID(ctx, *resp.MediaInfo.GetThumbnailID())
		if thumbnail != nil {
			if imageData := thumbnail.ImageData(); imageData != nil {
				return imageData, nil
			}
		}
	}

	if resp.YoutubeChannelID != nil {
		channel, _ := uc.FindYoutubeChannelInfo(ctx, *resp.YoutubeChannelID)
		if channel != nil && len(channel.ImageRaw) > 0 {
			return &dtypes.ImageData{
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
