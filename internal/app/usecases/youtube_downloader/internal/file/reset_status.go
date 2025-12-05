package fileuc

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *File) ResetStatus(ctx context.Context) error {
	statuses := []dtypes.FileStatus{
		dtypes.FileStatusPending,
		dtypes.FileStatusWorking,
	}

	err := uc.FileRep.UpdateStatusToNew(ctx, statuses)
	if err != nil {
		uc.logger.Warn("Failed update status to new", "error", err)
		return err
	}

	return nil
}
