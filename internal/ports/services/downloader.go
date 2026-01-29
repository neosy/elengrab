package pservices

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

type Downloader interface {
	GetTitle(ctx context.Context, url string) (string, error)
	GetFormats(ctx context.Context, url string) (*dmedia.MediaInfo, error)
	GetBestFormat(ctx context.Context, url string) (*dmedia.MediaInfo, error)
	Download(ctx context.Context, url string, options *dservices.DownloadOptions) (<-chan *ddownload.DownloadResult, error)
}
