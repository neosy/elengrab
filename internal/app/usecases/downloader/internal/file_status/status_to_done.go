package filestatus

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

// Failed set status to done
func (s *FileStatus) Done(
	ctx context.Context,
	fileId uuid.UUID,
	patch *dto.FileInfoPatch,
) error {
	updateFieldsFunc := func(file *ddownload.File) {
		dto.PatchToFileDomain(patch, file)
		file.DownloadedAt = uptr.Any(time.Now().UTC())
	}

	return s.file.Tx(ctx, func(ctx context.Context) error {
		err := s.dlTask.DeleteByFileId(ctx, fileId)
		if err != nil {
			return err
		}

		return s.updateStatus(
			ctx,
			fileId,
			dtypes.FileStatusDone,
			updateFieldsFunc,
		)
	})
}
