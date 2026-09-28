package downloader

import "context"

func (d *Downloader) GetVersion(ctx context.Context) (string, error) {
	return d.executor.GetVersion(ctx)
}
