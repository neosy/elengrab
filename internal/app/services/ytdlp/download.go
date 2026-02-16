package ytdlpsrv

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

func (srv *YtDlpService) Download(
	ctx context.Context,
	url string,
	options *dservices.DownloadOptions,
) (<-chan *ddownload.DownloadResult, error) {
	// Create a full channel (read/write)
	resultCh := make(chan *ddownload.DownloadResult)

	// Launch the goroutine that writes into the channel
	go func() {
		defer close(resultCh)
		srv.downloader.Download(ctx, url, options, resultCh)
	}()

	return resultCh, nil
}
