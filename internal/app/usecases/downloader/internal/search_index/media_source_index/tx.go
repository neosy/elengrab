package sourceindex

import "context"

func (uc *MediaSourceIndex) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.indexRepo().Tx(ctx, fn)
}

func (uc *MediaSourceIndex) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.indexRepo().TxIndependent(ctx, fn)
}
