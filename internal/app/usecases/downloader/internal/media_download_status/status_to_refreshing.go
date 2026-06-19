package downloadstatus

import (
	"context"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Refreshing set status to refreshing
func (s *MediaDownloadStatus) Refreshing(
	ctx context.Context,
	downloadID uuid.UUID,
) error {
	tx := func(ctx context.Context) error {
		err := s.updateStatus(
			ctx,
			downloadID,
			dtypes.MediaDownloadStatusRefreshing,
			nil,
		)
		if err != nil {
			return err
		}

		return nil
	}

	return s.download.Tx(ctx, tx)
}
