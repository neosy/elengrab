package ytdlpsrv

import (
	"context"

	pservices "github.com/neosy/elengrab/internal/ports/services"
)

// FetchInfoWithBestFormat retrieves and parses video best format for the given URL.
func (srv *YtDlpService) FetchInfoWithBestFormat(
	ctx context.Context,
	url string,
	opts ...pservices.DownloaderOption,
) (*pservices.DownloaderMeta, error) {
	bestInfo, err := srv.downloader.FetchInfoWithBestFormat(ctx, url, "b", internalRequestOptions(opts...))
	if err != nil {
		return nil, err
	}

	srv.logger.Debug(
		"Get best format",
		"url", url,
		"format", bestInfo,
	)

	return bestInfo, nil
}
