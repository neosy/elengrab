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

func (uc *Downloader) GetFileImage(
	ctx context.Context,
	userCtx dauth.UserContext,
	fileID uuid.UUID,
	sources []dtypes.ImageSource,
) (*dtypes.ImageData, error) {
	fileInfo, err := uc.GetFileInfo(ctx, userCtx, fileID)
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
			imageData, err = uc.getFileImageThumbnail(ctx, fileInfo)
		case dtypes.ImageSourceAvatar:
			imageData, err = uc.getFileImageAvatar(ctx, fileInfo)
		case dtypes.ImageSourceSite:
			imageData, err = uc.getFileImageSite(ctx, fileInfo)
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

func (uc *Downloader) getFileImageThumbnail(
	ctx context.Context,
	fileInfo *dto.GetFileInfoResponse,
) (*dtypes.ImageData, error) {
	if fileInfo.MediaInfo != nil && fileInfo.MediaInfo.GetThumbnailID() != nil {
		thumbnail, _ := uc.thumbnail.GetByThumbID(ctx, *fileInfo.MediaInfo.GetThumbnailID())
		if thumbnail != nil {
			if imageData := thumbnail.ImageData(); imageData != nil {
				return imageData, nil
			}
		}
	}

	return nil, errorx.NewHTTPMessage("thumbnail not found", http.StatusNotFound)
}

func (uc *Downloader) getFileImageAvatar(
	ctx context.Context,
	fileInfo *dto.GetFileInfoResponse,
) (*dtypes.ImageData, error) {
	if fileInfo.ChannelID != nil && fileInfo.IsYouTube() {
		channel, _ := uc.FindYoutubeChannelInfo(ctx, *fileInfo.ChannelID)
		if channel != nil && len(channel.ImageRaw) > 0 {
			return &dtypes.ImageData{
				Raw:    channel.ImageRaw,
				Format: channel.ImageFormat,
			}, nil
		}
	}

	return nil, errorx.NewHTTPMessage("avatar not found", http.StatusNotFound)
}

func (uc *Downloader) getFileImageSite(
	ctx context.Context,
	fileInfo *dto.GetFileInfoResponse,
) (*dtypes.ImageData, error) {
	logo, err := uc.siteIcon.GetBySiteURL(ctx, httpx.BaseURL(fileInfo.MediaUrl))
	if err != nil {
		return nil, errorx.Errorf("avatar not found: %w", err)
	}

	return logo.ImageData(), nil
}
