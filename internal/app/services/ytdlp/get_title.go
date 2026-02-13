package ytdlpsrv

import (
	"context"
)

// GetTitle
func (srv *YtDlpService) GetTitle(ctx context.Context, url string, useCookies bool) (string, error) {
	return srv.core.GetTitle(ctx, url, useCookies)
}
