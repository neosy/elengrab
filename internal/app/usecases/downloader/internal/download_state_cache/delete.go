package dlstate

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadStateCache) Delete(ctx context.Context, downloadID uuid.UUID) error {
	return uc.stateCacheRep.Delete(ctx, downloadID)
}
