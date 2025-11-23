package fileuc

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *File) DeleteBroken(ctx context.Context) error {
	files, err := uc.GetByStatus(ctx, dtypes.FileStatusNew)
	if err != nil {
		uc.logger.Warn("Failed get files", "error", err)
		return err
	}

	for _, file := range files {
		if file.DownloadTask == nil {
			err := uc.Delete(ctx, file.FileId)
			if err != nil {
				uc.logger.Warn("Failed delete file", "error", err)
			}
		}
	}

	return nil
}
