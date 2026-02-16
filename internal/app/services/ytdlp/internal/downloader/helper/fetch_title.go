package helper

import (
	"context"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	"github.com/neosy/elengrab/pkg/httpx"
)

func FetchTitleFast(ctx context.Context, url string) (string, error) {
	if isYouTube(url) {
		info, err := FetchYoutubeInfoFast(url)
		if err != nil {
			return "", err
		}
		return info.Title, nil
	}

	return httpx.GetTitle(
		ctx,
		url,
		httpx.ClientOptionWithDefaultCookieJar(),
		httpx.ClientOptionWithTimeout(consts.FetchTitleTimeout),
	)
}
