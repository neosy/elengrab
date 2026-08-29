package downloader

import "context"

func (uc *downloader) Start(ctx context.Context) {
	uc.mediaWatch.Start(ctx)
}
