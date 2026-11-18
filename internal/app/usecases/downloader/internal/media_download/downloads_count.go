package mediadownload

import (
	"context"
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *MediaDownload) IncompleteDownloadsCount(
	ctx context.Context,
	userID uuid.UUID,
) (dtypes.IncompleteDownloadsCount, error) {
	totalCount, err := uc.dlStateCache.Count(ctx, userID)
	if err != nil {
		return dtypes.IncompleteDownloadsCount{}, err
	}

	uploadsCount, err := uc.uploadStateCacheRepo.Count(ctx, userID)
	if err != nil {
		return dtypes.IncompleteDownloadsCount{}, err
	}

	totalCount.Add(uploadsCount)

	return totalCount, nil
}
