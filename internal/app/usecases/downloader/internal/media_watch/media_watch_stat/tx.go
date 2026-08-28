package watchstat

import "context"

func (uc *MediaWatchStat) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.statRepo().Tx(ctx, fn)
}

func (uc *MediaWatchStat) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.statRepo().TxIndependent(ctx, fn)
}
