package ytdlpsrv

import (
	"context"

	dservices "github.com/neosy/elengrab/internal/domain/services"
)

// FetchInfo retrieves and parses video formats for the given URL.
func (srv *YtDlpService) FetchInfo(
	ctx context.Context,
	url string,
	opts ...RequestOption,
) (*dservices.DownloaderMediaInfo, error) {
	return srv.downloader.FetchInfo(ctx, url, internalRequestOptions(opts...))
}
