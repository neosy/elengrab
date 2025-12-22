package ytdlpsrv

import (
	"context"

	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
)

// GetBestFormat retrieves and parses video best format for the given URL.
func (srv *YtDlpService) GetBestFormat(
	ctx context.Context,
	url string,
) (*dyoutubeinfo.YouTubeInfo, error) {
	return srv.getBestFormat(ctx, url, "b")
}

func (srv *YtDlpService) getBestFormat(
	ctx context.Context,
	url string,
	format string,
) (*dyoutubeinfo.YouTubeInfo, error) {
	bestInfo, err := srv.ytdlp.GetBestFormat(ctx, url, format)
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
