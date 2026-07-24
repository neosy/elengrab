package mediawatch

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaWatch) GetViews(ctx context.Context, downloadID uuid.UUID) (uint32, error) {
	stat, err := uc.stat.Find(ctx, downloadID)
	if err != nil {
		return 0, err
	}

	if stat == nil {
		return 0, nil
	}

	return stat.Views, nil
}
