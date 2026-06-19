package repository

import "context"

func (r *ThumbnailRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.repo.Tx(ctx, fn)
}
