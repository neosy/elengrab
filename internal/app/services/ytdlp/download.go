package ytdlpsrv

import (
	"context"

	dservices "github.com/neosy/elengrab/internal/domain/services"
)

// Download starts a download process for the given URL and returns a channel to receive download results.
func (srv *YtDlpService) Download(
	ctx context.Context,
	url string,
	options *dservices.DownloadOptions,
) (<-chan *dservices.DownloadResult, error) {
	// Create a full channel (read/write)
	resultCh := make(chan *dservices.DownloadResult)

	// Launch the goroutine that writes into the channel
	go func() {
		defer close(resultCh)
		srv.downloader.Download(ctx, url, options, resultCh)
	}()

	return resultCh, nil
}
