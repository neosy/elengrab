package uwatchchunk

import (
	"context"

	"uuid"
)

func (uc *MediaUserWatchChunk) IterateDownloadUsers(
	ctx context.Context,
	fn func(downloadID, userID uuid.UUID) error,
) error {
	return uc.chunkRepo().IterateDownloadUsers(ctx, fn)
}
