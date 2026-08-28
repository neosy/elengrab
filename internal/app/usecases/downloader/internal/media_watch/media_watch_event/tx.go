package watchevent

import "context"

func (uc *MediaWatchEvent) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.eventRepo().Tx(ctx, fn)
}

func (uc *MediaWatchEvent) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.eventRepo().TxIndependent(ctx, fn)
}
