package ytdlpsrv

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// GetBestFormat retrieves and parses video best format for the given URL.
func (srv *YtDlpService) GetBestFormat(
	ctx context.Context,
	url string,
	useCookies bool,
) (*dmedia.MediaInfo, error) {
	bestInfo, err := srv.downloader.GetBestFormat(ctx, url, "b", useCookies)
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
