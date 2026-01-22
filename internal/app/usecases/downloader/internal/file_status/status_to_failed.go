package filestatus

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// Failed set status to failed
func (s *FileStatus) Failed(
	ctx context.Context,
	fileId uuid.UUID,
	patch *dto.FileInfoPatch,
	message *string,
) error {
	updateFieldsFunc := func(file *ddownload.File) {
		dto.PatchToFileDomain(patch, file)
		file.ErrorMessage = message
	}

	err := s.dlTask.DeleteByFileId(ctx, fileId)
	if err != nil {
		return err
	}

	return s.updateStatus(
		ctx,
		fileId,
		dtypes.FileStatusFailed,
		updateFieldsFunc,
	)
}
