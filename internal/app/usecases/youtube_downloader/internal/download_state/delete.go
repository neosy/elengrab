package dlstate

import (
	"context"

	"github.com/google/uuid"
)

func (uc *DownloadState) Delete(ctx context.Context, fileId uuid.UUID) error {
	return uc.stateRep.Delete(ctx, fileId)
}
