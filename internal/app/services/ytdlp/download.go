package ytdlpsrv

import (
	"context"

	pservices "github.com/neosy/elengrab/internal/ports/services"
)

// Download starts a download process for the given URL and returns a channel to receive download results.
func (srv *YtDlpService) Download(
	ctx context.Context,
	url string,
	options *pservices.DownloadOptions,
) (<-chan *pservices.DownloadResult, error) {
	// Create a full channel (read/write)
	resultCh := make(chan *pservices.DownloadResult)

	// Launch the goroutine that writes into the channel
	go func() {
		defer close(resultCh)
		srv.downloader.Download(ctx, url, options, resultCh)
	}()

	return resultCh, nil
}
