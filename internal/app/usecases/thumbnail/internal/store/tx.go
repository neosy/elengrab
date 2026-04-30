package store

import "context"

func (uc *ThumbnailStore) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.thumbnailRep.Tx(ctx, fn)
}
