package filestatus

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Failed set status to done
func (s *FileStatus) Done(
	ctx context.Context,
	fileId uuid.UUID,
	patch *dto.FileInfoPatch,
) error {
	updateFieldsFunc := func(file *ddownload.File) {
		dto.PatchToFileDomain(patch, file)
	}

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
}
