package watchevent

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaWatchEvent) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	return uc.eventRep.DeleteByDownloadID(ctx, downloadID)
}
