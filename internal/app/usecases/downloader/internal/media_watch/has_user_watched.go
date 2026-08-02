package mediawatch

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaWatch) HasUserWatched(
	ctx context.Context,
	downloadID uuid.UUID, userID uuid.UUID,
) (bool, error) {
	return uc.userStat.Exists(ctx, downloadID, userID)
}
