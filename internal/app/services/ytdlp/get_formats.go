package ytdlpsrv

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// GetFormats retrieves and parses video formats for the given URL.
func (srv *YtDlpService) GetFormats(ctx context.Context, url string, useCookies bool) (*dmedia.MediaInfo, error) {
	return srv.downloader.GetFormats(ctx, url, useCookies)
}
