package mediadownload

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaDownload) UpdateUser(ctx context.Context, fromID, toID uuid.UUID) error {
	err := uc.downloadRep.UpdateOwner(ctx, fromID, toID)
	if err != nil {
		uc.logger.Warn("Update owner error", "fromID", fromID, "toID", toID, "error", err)
		return err
	}
	return nil
}
