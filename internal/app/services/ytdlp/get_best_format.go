package ytdlpsrv

import (
	"context"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// GetBestFormat retrieves and parses video best format for the given URL.
func (srv *YtDlpService) GetBestFormat(
	ctx context.Context,
	url string,
) (*dmedia.MediaInfo, error) {
	return srv.getBestFormat(ctx, url, "b")
}

func (srv *YtDlpService) getBestFormat(
	ctx context.Context,
	url string,
	format string,
) (*dmedia.MediaInfo, error) {
	bestInfo, err := srv.core.GetBestFormat(ctx, url, format)
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
