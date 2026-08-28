package sourceindex

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaSourceIndex) SoftDelete(ctx context.Context, downloadID uuid.UUID) error {
	err := uc.indexRepo().SoftDelete(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed soft delete MediaSourceIndex", "error", err)
		return err
	}

	return nil
}

func (uc *MediaSourceIndex) HardDelete(ctx context.Context, downloadID uuid.UUID) error {
	err := uc.indexRepo().HardDelete(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed hard delete MediaSourceIndex", "error", err)
		return err
	}

	return nil
}
