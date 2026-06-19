package ytdlpsrv

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (srv *YtDlpService) ExtractThumbnailURL(
	ctx context.Context,
	mediaURL string,
	opts ...RequestOption,
) (string, error) {
	return srv.downloader.ExtractThumbnailURL(ctx, mediaURL, internalFetchImageOptions(opts...))
}

func (srv *YtDlpService) FetchThumbnail(
	ctx context.Context,
	mediaURL string,
	opts ...RequestOption,
) (*dtypes.ImageData, error) {
	return srv.downloader.FetchThumbnail(ctx, mediaURL, internalFetchImageOptions(opts...))
}
