package downloader

import "context"

func (uc *Downloader) Start(ctx context.Context) {
	uc.mediaWatch.Start(ctx)
}
