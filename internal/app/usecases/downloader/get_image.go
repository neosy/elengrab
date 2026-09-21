package downloader

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

func (uc *downloader) GetDownloadImage(
	ctx context.Context,
	userCtx dauth.AuthContext,
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
			dtypes.ImageSourceChannel,
			dtypes.ImageSourceSite,
		}
	}

	var imageData *dtypes.ImageData

	for _, src := range sources {
		switch src {
		case dtypes.ImageSourceThumbnail:
			imageData, err = uc.getDownloadThumbnailImage(ctx, downloadInfo.MediaInfo)
		case dtypes.ImageSourceChannel:
			imageData, err = uc.GetChannelImage(ctx, downloadInfo.Channel)
		case dtypes.ImageSourceSite:
			imageData, err = uc.GetSiteImage(ctx, downloadInfo.MediaURL)
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

func (uc *downloader) getDownloadThumbnailImage(
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

func (uc *downloader) GetChannelImage(
	_ context.Context,
	channel *dmedia.Channel,
) (*dtypes.ImageData, error) {
	if channel != nil && channel.HasImage() {
		return channel.ImageData(), nil
	}

	return nil, errorx.NewHTTPMessage("avatar not found", http.StatusNotFound)
}

func (uc *downloader) GetSiteImage(
	ctx context.Context,
	url string,
) (*dtypes.ImageData, error) {
	logo, err := uc.siteIcon.GetBySiteURL(ctx, httpx.BaseURL(url))
	if err != nil {
		return nil, errorx.Errorf("avatar not found: %w", err)
	}

	return logo.ImageData(), nil
}
