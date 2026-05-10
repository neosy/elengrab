package pservices

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type Downloader interface {
	GetTitle(ctx context.Context, url string, useCookies bool) (string, error)

	GetInfo(ctx context.Context, url string, useCookies bool) (*dmedia.MediaInfo, error)
	GetInfoWithBestFormat(ctx context.Context, url string, useCookies bool) (*dmedia.MediaInfo, error)

	Download(ctx context.Context, url string, options *dservices.DownloadOptions) (<-chan *dservices.DownloadResult, error)

	ExtractThumbnailURL(ctx context.Context, mediaURL string) (string, error)
	FetchThumbnail(ctx context.Context, mediaURL string) (*dtypes.ImageData, error)
}
