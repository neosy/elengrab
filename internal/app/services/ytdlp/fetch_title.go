package ytdlpsrv

import (
	"context"

	pservices "github.com/neosy/elengrab/internal/ports/services"
)

// FetchTitle retrieves title for the given URL.
func (srv *YtDlpService) FetchTitle(ctx context.Context, url string, opts ...pservices.DownloaderOption) (string, error) {
	return srv.downloader.FetchTitle(ctx, url, internalRequestOptions(opts...))
}
