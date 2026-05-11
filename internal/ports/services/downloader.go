package pservices

import (
	"context"

	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type Downloader interface {
	GetTitle(ctx context.Context, url string, useCookies bool) (string, error)

	GetInfo(ctx context.Context, url string, useCookies bool) (*dservices.DownloaderMediaInfo, error)
	GetInfoWithBestFormat(ctx context.Context, url string, useCookies bool) (*dservices.DownloaderMediaInfo, error)

	Download(ctx context.Context, url string, options *dservices.DownloadOptions) (<-chan *dservices.DownloaderResult, error)

	ExtractThumbnailURL(ctx context.Context, mediaURL string, useCookies bool) (string, error)
	FetchThumbnail(ctx context.Context, mediaURL string, useCookies bool) (*dtypes.ImageData, error)
}
