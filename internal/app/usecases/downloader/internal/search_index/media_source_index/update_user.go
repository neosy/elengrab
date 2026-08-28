package sourceindex

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaSourceIndex) UpdateUser(ctx context.Context, fromID, toID uuid.UUID) error {
	err := uc.indexRepo().UpdateOwner(ctx, fromID, toID)
	if err != nil {
		uc.logger.Warn(
			"Update owner error in MediaSourceIndex",
			"fromID", fromID,
			"toID", toID,
			"error", err,
		)
		return err
	}

	return nil
}
