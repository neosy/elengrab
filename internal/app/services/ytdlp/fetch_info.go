package ytdlpsrv

import (
	"context"

	pservices "github.com/neosy/elengrab/internal/ports/services"
)

// FetchInfo retrieves and parses video formats for the given URL.
func (srv *YtDlpService) FetchInfo(
	ctx context.Context,
	url string,
	opts ...pservices.DownloaderOption,
) (*pservices.DownloaderMeta, error) {
	return srv.downloader.FetchInfo(ctx, url, internalRequestOptions(opts...))
}
