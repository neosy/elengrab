package fileuc

import (
	"context"

	"github.com/google/uuid"
)

func (uc *File) UpdateUser(ctx context.Context, fromID, toID uuid.UUID) error {
	err := uc.fileRep.UpdateOwner(ctx, fromID, toID)
	if err != nil {
		uc.logger.Warn("Update owner error", "fromID", fromID, "toID", toID, "error", err)
		return err
	}
	return nil
}
