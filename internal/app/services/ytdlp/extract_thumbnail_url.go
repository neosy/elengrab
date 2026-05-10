package ytdlpsrv

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (srv *YtDlpService) ExtractThumbnailURL(ctx context.Context, mediaURL string) (string, error) {
	return srv.downloader.ExtractThumbnailURL(ctx, mediaURL)
}

func (srv *YtDlpService) FetchThumbnail(ctx context.Context, mediaURL string) (*dtypes.ImageData, error) {
	return srv.downloader.FetchThumbnail(ctx, mediaURL)
}
