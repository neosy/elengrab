package searchindex

import (
	"context"

	"github.com/google/uuid"
)

func (uc *SearchIndex) UpdateViews(ctx context.Context, downloadID uuid.UUID, views uint32) error {
	return uc.searchIndex.UpdateViews(ctx, downloadID, views)
}
