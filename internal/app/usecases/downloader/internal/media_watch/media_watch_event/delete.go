package watchevent

import (
	"context"

	"uuid"
)

func (uc *MediaWatchEvent) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	return uc.eventRepo().DeleteByDownloadID(ctx, downloadID)
}
