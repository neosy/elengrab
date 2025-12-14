package ytdlpsrv

import (
	"context"
)

// GetTitle
func (srv *YtDlpService) GetTitle(ctx context.Context, url string) (string, error) {
	return srv.ytdlp.GetTitle(ctx, url)
}
