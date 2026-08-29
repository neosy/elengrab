package downloader

import (
	"context"

	"github.com/google/uuid"
)

func (uc *downloader) MigrateGuestData(ctx context.Context, guestID, userID uuid.UUID) error {
	return uc.download.UpdateUser(ctx, guestID, userID)
}
