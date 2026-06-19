package repository

import (
	"context"

	"github.com/google/uuid"
)

func (r *ThumbnailRepository) Delete(ctx context.Context, thumbID uuid.UUID) error {
	err := r.repo.Delete(ctx, thumbID)
	if err != nil {
		r.logger.Warn("Failed to delete thumbnail",
			"thumbnailID", thumbID,
			"error", err,
		)
		return err
	}

	return nil
}
