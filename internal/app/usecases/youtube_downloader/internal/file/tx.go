package fileuc

import "context"

func (uc *File) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.fileRep.Tx(ctx, fn)
}
