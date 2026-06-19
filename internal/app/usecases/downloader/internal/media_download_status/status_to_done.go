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
	patch func(*ddownload.MediaDownload),
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
			func(download *ddownload.MediaDownload) {
				if patch != nil {
					patch(download)
				}
				if download.DownloadedAt == nil {
					download.DownloadedAt = new(time.Now().UTC())
				}
			},
		)
		if err != nil {
			return err
		}

		return nil
	}

	return s.download.Tx(ctx, update)
}
