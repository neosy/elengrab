package watchstat

import "context"

func (uc *MediaWatchStat) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.statRep.Tx(ctx, fn)
}

func (uc *MediaWatchStat) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.statRep.TxIndependent(ctx, fn)
}
