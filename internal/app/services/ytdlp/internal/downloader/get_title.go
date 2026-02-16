package downloader

import "context"

func (d *Downloader) GetTitle(ctx context.Context, url string, useCookies bool) (string, error) {
	return d.executor.GetTitle(ctx, url, useCookies)
}
