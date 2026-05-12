package dlstate

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadStateCache) Delete(ctx context.Context, downloadID uuid.UUID) error {
	return uc.stateRep.Delete(ctx, downloadID)
}
