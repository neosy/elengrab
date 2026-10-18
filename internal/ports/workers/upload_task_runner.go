package pworkers

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type UploadTaskRunner interface {
	ExecuteUploadTask(ctx context.Context, workerID uint64, task *ddownload.MediaProcessingTask) error
}
