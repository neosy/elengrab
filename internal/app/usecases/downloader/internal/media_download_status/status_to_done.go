package downloadstatus

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Failed set status to done
func (s *MediaDownloadStatus) Done(
	ctx context.Context,
	downloadID uuid.UUID,
	mutate func(*ddownload.MediaDownload) error,
) error {
	update := func(ctx context.Context) error {
		err := s.dlTask.DeleteByDownloadID(ctx, downloadID)
		if err != nil {
			return err
		}

		err = s.updateStatus(
			ctx,
			downloadID,
			dtypes.MediaDownloadStatusDone,
			func(download *ddownload.MediaDownload) error {
				if mutate != nil {
					if err := mutate(download); err != nil {
						return err
					}
				}
				if download.DownloadedAt == nil {
					download.DownloadedAt = new(time.Now().UTC())
				}
				return nil
			},
		)
		if err != nil {
			return err
		}

		return nil
	}

	return s.download.Tx(ctx, update)
}
