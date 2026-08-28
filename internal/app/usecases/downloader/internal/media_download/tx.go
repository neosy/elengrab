package mediadownload

import "context"

func (uc *MediaDownload) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.downloadRep().Tx(ctx, fn)
}
