package thumbnails

import (
	"context"

	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func FetchImages(
	ctx context.Context,
	mediaURL string,
	options FetchOptions,
) ([]*dtypes.ImageData, error) {
	switch hostdetect.DetectPlatformType(mediaURL) {
	case dtypes.MediaPlatformTypeYouTube:
		return FetchYouTubeImages(ctx, mediaURL, options)
	}

	return nil, nil
}
