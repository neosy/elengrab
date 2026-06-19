package downloader

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

func (uc *Downloader) GetDownloadImage(
	ctx context.Context,
	userCtx dauth.UserContext,
	downloadID uuid.UUID,
	sources []dtypes.ImageSource,
) (*dtypes.ImageData, error) {
	downloadInfo, err := uc.GetDownloadInfo(ctx, userCtx, downloadID)
	if err != nil {
		return nil, err
	}

	if len(sources) == 0 {
		sources = []dtypes.ImageSource{
			dtypes.ImageSourceThumbnail,
			dtypes.ImageSourceAvatar,
			dtypes.ImageSourceSite,
		}
	}

	var imageData *dtypes.ImageData

	for _, src := range sources {
		switch src {
		case dtypes.ImageSourceThumbnail:
			imageData, err = uc.getDownloadImageThumbnail(ctx, downloadInfo.MediaInfo)
		case dtypes.ImageSourceAvatar:
			imageData, err = uc.getDownloadImageAvatar(ctx, downloadInfo)
		case dtypes.ImageSourceSite:
			imageData, err = uc.getDownloadImageSite(ctx, downloadInfo)
		}
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, errorx.Errorf("image not found: %w", err)
	}

	if imageData == nil {
		return nil, errorx.NewHTTPMessage("image not found", http.StatusNotFound)
	}

	return imageData, nil
}

func (uc *Downloader) getDownloadImageThumbnail(
	ctx context.Context,
	mediaInfo *dtypes.MediaInfo,
) (*dtypes.ImageData, error) {
	if mediaInfo != nil && mediaInfo.PreferredThumbnailID() != uuid.Nil {
		thumbnail, _ := uc.thumbnail.LoadByThumbID(ctx, mediaInfo.PreferredThumbnailID())
		if thumbnail != nil {
			if imageData := thumbnail.ImageData(""); imageData != nil {
				return imageData, nil
			}
		}
	}

	return nil, errorx.NewHTTPMessage("thumbnail not found", http.StatusNotFound)
}

func (uc *Downloader) getDownloadImageAvatar(
	ctx context.Context,
	downloadInfo *dto.GetMediaDownloadInfoResponse,
) (*dtypes.ImageData, error) {
	if downloadInfo.ChannelID != nil && downloadInfo.IsYouTube() {
		channel, _ := uc.FindYoutubeChannelInfo(ctx, *downloadInfo.ChannelID)
		if channel != nil && len(channel.ImageRaw) > 0 {
			return &dtypes.ImageData{
				Raw:    channel.ImageRaw,
				Format: channel.ImageFormat,
			}, nil
		}
	}

	return nil, errorx.NewHTTPMessage("avatar not found", http.StatusNotFound)
}

func (uc *Downloader) getDownloadImageSite(
	ctx context.Context,
	downloadInfo *dto.GetMediaDownloadInfoResponse,
) (*dtypes.ImageData, error) {
	logo, err := uc.siteIcon.GetBySiteURL(ctx, httpx.BaseURL(downloadInfo.MediaURL))
	if err != nil {
		return nil, errorx.Errorf("avatar not found: %w", err)
	}

	return logo.ImageData(), nil
}
