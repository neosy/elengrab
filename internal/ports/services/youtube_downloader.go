package pservices

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
)

type YouTubeDownloader interface {
	GetTitle(ctx context.Context, url string) (string, error)
	GetFormats(ctx context.Context, url string) (*dyoutubeinfo.YouTubeInfo, error)
	GetBestFormat(ctx context.Context, url string) (*dyoutubeinfo.YouTubeInfo, error)
	Download(ctx context.Context, url string, options *dservices.DownloadOptions) (<-chan *ddownload.DownloadResult, error)
}
