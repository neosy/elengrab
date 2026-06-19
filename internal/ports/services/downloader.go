package pservices

import (
	"context"

	ytdlpsrv "github.com/neosy/elengrab/internal/app/services/ytdlp"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type Downloader interface {
	FetchTitle(ctx context.Context, url string, opts ...ytdlpsrv.RequestOption) (string, error)

	FetchInfo(ctx context.Context, url string, opts ...ytdlpsrv.RequestOption) (*dservices.DownloaderMediaInfo, error)
	FetchInfoWithBestFormat(ctx context.Context, url string, opts ...ytdlpsrv.RequestOption) (*dservices.DownloaderMediaInfo, error)

	Download(ctx context.Context, url string, options *dservices.DownloadOptions) (<-chan *dservices.DownloaderResult, error)

	ExtractThumbnailURL(ctx context.Context, mediaURL string, opts ...ytdlpsrv.RequestOption) (string, error)
	FetchThumbnail(ctx context.Context, mediaURL string, opts ...ytdlpsrv.RequestOption) (*dtypes.ImageData, error)
}
