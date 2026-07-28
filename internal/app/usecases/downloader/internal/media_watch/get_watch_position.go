package mediawatch

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (uc *MediaWatch) GetLastUserWatchPosition(
	ctx context.Context,
	downloadID uuid.UUID, userID uuid.UUID, sessionID *uuid.UUID,
) (time.Duration, error) {
	position, err := uc.userPosition.Find(ctx, downloadID, userID, sessionID)
	if err != nil {
		return 0, err
	}

	if position == nil {
		return 0, nil
	}

	return position.Position, nil
}
