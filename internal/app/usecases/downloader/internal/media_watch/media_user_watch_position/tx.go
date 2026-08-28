package uwatchposition

import "context"

func (uc *MediaUserWatchPosition) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.positionRepo().Tx(ctx, fn)
}

func (uc *MediaUserWatchPosition) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.positionRepo().TxIndependent(ctx, fn)
}
