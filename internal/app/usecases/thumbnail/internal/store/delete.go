package store

import (
	"context"

	"github.com/google/uuid"
)

func (c *ThumbnailStore) Delete(ctx context.Context, thumbID uuid.UUID) error {
	err := c.thumbnailRep.Delete(ctx, thumbID)
	if err != nil {
		c.logger.Warn("Failed to delete thumbnail",
			"thumbnailID", thumbID,
			"error", err,
		)
		return err
	}

	return nil
}
