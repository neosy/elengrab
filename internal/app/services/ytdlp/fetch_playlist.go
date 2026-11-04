package ytdlpsrv

import (
	"context"

	pservices "github.com/neosy/elengrab/internal/ports/services"
)

// FetchPlaylist fetches a playlist from the provided URL and returns a DownloaderPlaylist object.
func (srv *YtDlpService) FetchPlaylist(
	ctx context.Context,
	playlistURL string,
) (*pservices.DownloaderPlaylist, error) {
	return srv.downloader.FetchPlaylistWithCookieFallback(ctx, playlistURL)
}
