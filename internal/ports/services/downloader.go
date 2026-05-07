package pservices

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

type Downloader interface {
	GetTitle(ctx context.Context, url string, useCookies bool) (string, error)
	GetInfo(ctx context.Context, url string, useCookies bool) (*dmedia.MediaInfo, error)
	GetInfoWithBestFormat(ctx context.Context, url string, useCookies bool) (*dmedia.MediaInfo, error)
	Download(ctx context.Context, url string, options *dservices.DownloadOptions) (<-chan *dservices.DownloadResult, error)
}
