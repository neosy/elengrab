package downloader

import (
	"context"
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *downloader) GetIncompleteDownloadsCount(
	ctx context.Context,
	userID uuid.UUID,
) (dtypes.IncompleteDownloadsCount, error) {
	return uc.download.IncompleteDownloadsCount(ctx, userID)
}
