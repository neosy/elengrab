package ytdlpsrv

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// FetchChannelInfo retrieves channel info for the given media URL.
func (srv *YtDlpService) FetchChannelInfoWithCookieFallback(
	ctx context.Context,
	mediaURL string,
) (*dtypes.ChannelSource, error) {
	return srv.downloader.FetchChannelInfoWithCookieFallback(ctx, mediaURL)
}
