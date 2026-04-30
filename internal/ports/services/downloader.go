package pservices

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

type Downloader interface {
	GetTitle(ctx context.Context, url string, useCookies bool) (string, error)
	GetFormats(ctx context.Context, url string, useCookies bool) (*dmedia.MediaInfo, error)
	GetBestFormat(ctx context.Context, url string, useCookies bool) (*dmedia.MediaInfo, error)
	Download(ctx context.Context, url string, options *dservices.DownloadOptions) (<-chan *dservices.DownloadResult, error)
}
