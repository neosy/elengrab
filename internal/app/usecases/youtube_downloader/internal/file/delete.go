package fileuc

import (
	"context"

	"github.com/google/uuid"
)

func (uc *File) SoftDelete(ctx context.Context, fileId uuid.UUID) error {
	err := uc.FileRep.Delete(ctx, fileId, true)
	if err != nil {
		uc.logger.Warn("Failed delete file", "error", err)
		return err
	}
	return nil
}

func (uc *File) HardDelete(ctx context.Context, fileId uuid.UUID) error {
	err := uc.FileRep.Delete(ctx, fileId, false)
	if err != nil {
		uc.logger.Warn("Failed delete file", "error", err)
		return err
	}
	return nil
}
