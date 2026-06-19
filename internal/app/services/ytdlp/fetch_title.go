package ytdlpsrv

import (
	"context"
)

// FetchTitle retrieves title for the given URL.
func (srv *YtDlpService) FetchTitle(ctx context.Context, url string, opts ...RequestOption) (string, error) {
	return srv.downloader.FetchTitle(ctx, url, internalRequestOptions(opts...))
}
