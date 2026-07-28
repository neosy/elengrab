package watchchunk

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaWatchChunk) IterateDownloadUsers(
	ctx context.Context,
	fn func(downloadID, userID uuid.UUID) error,
) error {
	return uc.chunkRep.IterateDownloadUsers(ctx, fn)
}
