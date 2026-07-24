package usecases

import "context"

func (uc *Usecases) Start(ctx context.Context) {
	uc.Downloader.Start(ctx)
}
