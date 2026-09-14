package thumbnails

import (
	"context"

	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func FetchImages(
	ctx context.Context,
	mediaURL string,
	opts ...FetchOption,
) ([]*dtypes.ImageData, error) {
	switch hostdetect.DetectPlatformType(mediaURL) {
	case dtypes.MediaPlatformTypeYouTube:
		return FetchYouTubeImages(ctx, mediaURL, opts...)
	}

	return nil, nil
}
