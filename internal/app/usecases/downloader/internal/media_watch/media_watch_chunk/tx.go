package watchchunk

import "context"

func (uc *MediaWatchChunk) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.chunkRep.Tx(ctx, fn)
}
