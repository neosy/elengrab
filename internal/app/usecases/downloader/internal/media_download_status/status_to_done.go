package downloadstatus

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

// Failed set status to done
func (s *MediaDownloadStatus) Done(
	ctx context.Context,
	downloadID uuid.UUID,
	patch *dto.MediaDownloadInfoPatch,
) error {
	updateFieldsFunc := func(file *ddownload.MediaDownload) {
		dto.PatchToMediaDownloadDomain(patch, file)
		file.DownloadedAt = uptr.Any(time.Now().UTC())
	}

	return s.download.Tx(ctx, func(ctx context.Context) error {
		err := s.dlTask.DeleteByDownloadID(ctx, downloadID)
		if err != nil {
			return err
		}

		return s.updateStatus(
			ctx,
			downloadID,
			dtypes.MediaDownloadStatusDone,
			updateFieldsFunc,
		)
	})
}
