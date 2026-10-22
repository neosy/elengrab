package searchindex

import (
	"context"

	"uuid"
)

func (uc *SearchIndex) UpdateUser(ctx context.Context, fromID, toID uuid.UUID) error {
	err := uc.sourceIndex.UpdateUser(ctx, fromID, toID)
	if err != nil {
		return err
	}
	return nil
}
