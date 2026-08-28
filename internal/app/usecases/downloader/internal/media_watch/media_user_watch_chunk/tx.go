package uwatchchunk

import "context"

func (uc *MediaUserWatchChunk) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.chunkRepo().Tx(ctx, fn)
}
