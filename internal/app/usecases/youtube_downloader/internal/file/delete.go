package fileuc

import (
	"context"

	"github.com/google/uuid"
)

func (uc *File) Delete(ctx context.Context, fileId uuid.UUID) error {
	err := uc.FileRep.Delete(ctx, fileId)
	if err != nil {
		uc.logger.Warn("Failed delete file", "error", err)
	}
	return nil
}
