package mediawatch

import (
	"context"

	"uuid"
)

func (uc *MediaWatch) HasUserWatched(
	ctx context.Context,
	downloadID uuid.UUID, userID uuid.UUID,
) (bool, error) {
	return uc.userStat.Exists(ctx, downloadID, userID)
}
