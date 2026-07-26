package watchposition

import "context"

func (uc *MediaWatchPosition) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.positionRep.Tx(ctx, fn)
}

func (uc *MediaWatchPosition) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.positionRep.TxIndependent(ctx, fn)
}
