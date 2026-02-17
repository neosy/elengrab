package downloader

import (
	"context"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/executor"
)

func (d *Downloader) GetTitle(ctx context.Context, url string, useCookies bool) (string, error) {
	return d.executor.GetTitle(ctx, url, executor.WithUseCookies(useCookies))
}
