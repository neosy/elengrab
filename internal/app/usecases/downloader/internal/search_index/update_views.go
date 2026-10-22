package searchindex

import (
	"context"

	"uuid"
)

func (uc *SearchIndex) UpdateViews(ctx context.Context, downloadID uuid.UUID, views uint32) error {
	return uc.sourceIndex.UpdateViews(ctx, downloadID, views)
}
