package mediadownload

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaDownload) UpdateUser(ctx context.Context, fromID, toID uuid.UUID) error {
	return uc.downloadRepo().Tx(ctx, func(ctx context.Context) error {
		err := uc.downloadRepo().UpdateOwner(ctx, fromID, toID)
		if err != nil {
			uc.logger.Warn("Update owner error", "fromID", fromID, "toID", toID, "error", err)
			return err
		}

		err = uc.searchIndex.UpdateUser(ctx, fromID, toID)
		if err != nil {
			return err
		}

		return nil
	})

}

func (uc *MediaDownload) UpdateChannelID(ctx context.Context, oldChannelID string, newChannelID uuid.UUID) error {
	return uc.downloadRepo().UpdateChannelID(ctx, oldChannelID, newChannelID)
}
