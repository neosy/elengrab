package downloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/images"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (h *DownloaderHandlers) thumbnailURLWithFallback(mediaInfo *dtypes.MediaInfo) string {
	if mediaInfo == nil || mediaInfo.FormatType == dtypes.FormatTypeNone {
		return images.ThumbnailDefault().URL
	}

	if thumbID := mediaInfo.PreferredThumbnailID(); thumbID != uuid.Nil {
		return httppaths.BuildThumbnailPath(thumbID)
	}

	if mediaInfo.FormatType == dtypes.FormatTypeAudioOnly {
		return images.ThumbnailMusicDefault().URL
	}

	return images.ThumbnailVideoDefault().URL
}

func (h *DownloaderHandlers) thumbnailImageData(ctx context.Context, mediaInfo *dtypes.MediaInfo) *dtypes.ImageData {
	var imageData *dtypes.ImageData

	if mediaInfo == nil || mediaInfo.PreferredThumbnailID() == uuid.Nil {
		return nil
	}

	thumbnail, _ := h.thumbnail.GetByThumbID(ctx, mediaInfo.PreferredThumbnailID())
	if thumbnail != nil {
		imageData = thumbnail.ImageData(httppaths.BuildThumbnailPath(thumbnail.ThumbID))

	}

	return imageData
}

func (h *DownloaderHandlers) thumbnailImageDataWithFallback(ctx context.Context, mediaInfo *dtypes.MediaInfo) *dtypes.ImageData {
	imageData := h.thumbnailImageData(ctx, mediaInfo)
	if imageData != nil {
		return imageData
	}

	if mediaInfo == nil || mediaInfo.FormatType == dtypes.FormatTypeNone {
		return images.ThumbnailDefault().ImageData()
	}

	if mediaInfo.FormatType == dtypes.FormatTypeAudioOnly {
		return images.ThumbnailMusicDefault().ImageData()
	}

	return images.ThumbnailVideoDefault().ImageData()
}
