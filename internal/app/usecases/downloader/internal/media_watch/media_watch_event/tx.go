package watchevent

import "context"

func (uc *MediaWatchEvent) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.eventRep.Tx(ctx, fn)
}

func (uc *MediaWatchEvent) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.eventRep.TxIndependent(ctx, fn)
}
