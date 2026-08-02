package mediawatch

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaWatch) DeleteAllByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	deleteActions := func(ctx context.Context) error {
		err := uc.userStat.DeleteByDownloadID(ctx, downloadID)
		if err != nil {
			return err
		}

		err = uc.stat.DeleteByDownloadID(ctx, downloadID)
		if err != nil {
			return err
		}

		err = uc.userChunk.DeleteByDownloadID(ctx, downloadID)
		if err != nil {
			return err
		}

		err = uc.event.DeleteByDownloadID(ctx, downloadID)
		if err != nil {
			return err
		}

		err = uc.userPosition.DeleteByDownloadID(ctx, downloadID)
		if err != nil {
			return err
		}

		return nil
	}

	return uc.event.TxIndependent(ctx, deleteActions)
}
