package watchevent

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaWatchEvent) Delete(ctx context.Context, downloadID uuid.UUID) error {
	return uc.eventRep.Delete(ctx, downloadID)
}
