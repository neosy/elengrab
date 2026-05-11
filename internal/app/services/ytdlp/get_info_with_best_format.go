package ytdlpsrv

import (
	"context"

	dservices "github.com/neosy/elengrab/internal/domain/services"
)

// GetBestFormat retrieves and parses video best format for the given URL.
func (srv *YtDlpService) GetInfoWithBestFormat(
	ctx context.Context,
	url string,
	useCookies bool,
) (*dservices.DownloaderMediaInfo, error) {
	bestInfo, err := srv.downloader.GetInfoWithBestFormat(ctx, url, "b", useCookies)
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
