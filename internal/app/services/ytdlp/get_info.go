package ytdlpsrv

import (
	"context"

	dservices "github.com/neosy/elengrab/internal/domain/services"
)

// GetInfo retrieves and parses video formats for the given URL.
func (srv *YtDlpService) GetInfo(ctx context.Context, url string, useCookies bool) (*dservices.DownloaderMediaInfo, error) {
	return srv.downloader.GetInfo(ctx, url, useCookies)
}
