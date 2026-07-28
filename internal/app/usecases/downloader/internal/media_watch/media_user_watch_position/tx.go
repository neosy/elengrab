package uwatchposition

import "context"

func (uc *MediaUserWatchPosition) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.positionRep.Tx(ctx, fn)
}

func (uc *MediaUserWatchPosition) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.positionRep.TxIndependent(ctx, fn)
}
