package helper

import (
	"context"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

func FetchTitleFast(ctx context.Context, url string) (string, error) {
	if hostdetect.YouTube(url) {
		info, err := FetchYoutubeInfoFast(url)
		if err != nil {
			return "", err
		}
		return info.Title, nil
	}

	return httpx.FetchTitle(
		ctx,
		url,
		httpx.ClientOptionWithDefaultCookieJar(),
		httpx.ClientOptionWithTimeout(consts.FetchTitleTimeout),
	)
}
