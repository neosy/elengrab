package downloadstatus

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

// Failed set status to done
func (s *MediaDownloadStatus) Done(
	ctx context.Context,
	downloadID uuid.UUID,
	patch func(*ddownload.MediaDownload),
) error {
	tx := func(ctx context.Context) error {
		err := s.dlTask.DeleteByDownloadID(ctx, downloadID)
		if err != nil {
			return err
		}

		err = s.updateStatus(
			ctx,
			downloadID,
			dtypes.MediaDownloadStatusDone,
			func(download *ddownload.MediaDownload) {
				patch(download)
				download.DownloadedAt = uptr.Any(time.Now().UTC())
			},
		)
		if err != nil {
			return err
		}

		return nil
	}

	return s.download.Tx(ctx, tx)
}
