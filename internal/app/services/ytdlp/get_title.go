package ytdlpsrv

import (
	"context"
)

// GetTitle retrieves title for the given URL.
func (srv *YtDlpService) GetTitle(ctx context.Context, url string, useCookies bool) (string, error) {
	return srv.downloader.GetTitle(ctx, url, useCookies)
}
