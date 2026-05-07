package ytdlpsrv

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// GetInfo retrieves and parses video formats for the given URL.
func (srv *YtDlpService) GetInfo(ctx context.Context, url string, useCookies bool) (*dmedia.MediaInfo, error) {
	return srv.downloader.GetInfo(ctx, url, useCookies)
}
