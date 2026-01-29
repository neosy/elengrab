package ytdlpsrv

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// GetFormats retrieves and parses video formats for the given URL.
func (srv *YtDlpService) GetFormats(ctx context.Context, url string) (*dmedia.MediaInfo, error) {
	return srv.core.GetFormats(ctx, url)
}
