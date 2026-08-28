package uwatchstat

import "context"

func (uc *MediaUserWatchStat) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.statRepo().Tx(ctx, fn)
}

func (uc *MediaUserWatchStat) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.statRepo().TxIndependent(ctx, fn)
}
