package fileuc

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *File) DeleteBroken(ctx context.Context) error {
	files, err := uc.GetByStatus(ctx, dtypes.FileStatusNew)
	if err != nil {
		uc.logger.Warn("Failed to get files", "error", err)
		return err
	}

	for _, file := range files {
		if file.DownloadTask == nil {
			err := uc.HardDelete(ctx, file.FileID)
			if err != nil {
				uc.logger.Warn("Failed to delete file", "fileId", file.FileID, "error", err)
				continue
			}
		}
	}

	return nil
}
