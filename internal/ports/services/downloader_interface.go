package pservices

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type Downloader interface {
	GetVersion(ctx context.Context) (string, error)

	FetchTitle(ctx context.Context, url string, opts ...DownloaderOption) (string, error)

	FetchInfo(ctx context.Context, url string, opts ...DownloaderOption) (*DownloaderMeta, error)
	FetchInfoWithBestFormat(ctx context.Context, url string, opts ...DownloaderOption) (*DownloaderMeta, error)

	FetchChannelInfoWithCookieFallback(ctx context.Context, mediaURL string) (*dtypes.ChannelSource, error)

	Download(ctx context.Context, url string, options *DownloadOptions) (<-chan *DownloadResult, error)

	ExtractThumbnailURL(ctx context.Context, mediaURL string, opts ...DownloaderOption) (string, error)
	FetchThumbnail(ctx context.Context, mediaURL string, opts ...DownloaderOption) (*dtypes.ImageData, error)
}
