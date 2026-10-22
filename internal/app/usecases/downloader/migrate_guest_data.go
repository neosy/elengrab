package downloader

import (
	"context"

	"uuid"
)

func (uc *downloader) MigrateGuestData(ctx context.Context, guestID, userID uuid.UUID) error {
	return uc.download.UpdateUser(ctx, guestID, userID)
}
