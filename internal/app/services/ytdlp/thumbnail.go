package ytdlpsrv

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func (srv *YtDlpService) ExtractThumbnailURL(
	ctx context.Context,
	mediaURL string,
	opts ...pservices.DownloaderOption,
) (string, error) {
	return srv.downloader.ExtractThumbnailURL(ctx, mediaURL, internalFetchImageOptions(opts...))
}

func (srv *YtDlpService) FetchThumbnail(
	ctx context.Context,
	mediaURL string,
	opts ...pservices.DownloaderOption,
) (*dtypes.ImageData, error) {
	return srv.downloader.FetchThumbnail(ctx, mediaURL, internalFetchImageOptions(opts...))
}
