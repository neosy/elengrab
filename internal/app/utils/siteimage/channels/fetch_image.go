package channels

import (
	"context"

	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func FetchImages(
	ctx context.Context,
	channelURL string,
	opts ...FetchOption,
) ([]dtypes.ImageData, error) {
	switch hostdetect.DetectPlatformType(channelURL) {
	case dtypes.MediaPlatformTypeYouTube:
		return FetchYouTubeImages(ctx, channelURL, opts...)
	}

	return nil, nil
}
