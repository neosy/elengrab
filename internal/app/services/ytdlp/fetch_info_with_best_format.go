package ytdlpsrv

import (
	"context"

	dservices "github.com/neosy/elengrab/internal/domain/services"
)

// FetchInfoWithBestFormat retrieves and parses video best format for the given URL.
func (srv *YtDlpService) FetchInfoWithBestFormat(
	ctx context.Context,
	url string,
	opts ...RequestOption,
) (*dservices.DownloaderMediaInfo, error) {
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
