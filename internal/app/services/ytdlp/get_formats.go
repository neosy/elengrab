package ytdlpsrv

import (
	"context"

	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
)

// GetFormats retrieves and parses video formats for the given URL.
func (srv *YtDlpService) GetFormats(ctx context.Context, url string) (*dyoutubeinfo.YouTubeInfo, error) {
	return srv.ytdlp.GetFormats(ctx, url)
}
