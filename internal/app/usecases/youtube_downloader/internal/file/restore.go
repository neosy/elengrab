package fileuc

import (
	"context"

	"github.com/google/uuid"
)

func (uc *File) Restore(ctx context.Context, fileId uuid.UUID) error {
	err := uc.FileRep.Restore(ctx, fileId)
	if err != nil {
		uc.logger.Warn("Failed restore file", "error", err)
		return err
	}
	return nil
}
