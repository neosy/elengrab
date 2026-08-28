package dltask

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *DownloadTask) DeleteByStatus(ctx context.Context, status dtypes.DownloadTaskStatus) error {
	return uc.TaskRepo().DeleteByStatus(ctx, status)
}
