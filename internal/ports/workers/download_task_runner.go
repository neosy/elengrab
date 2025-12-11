package pworkers

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type DownloadTaskRunner interface {
	ExecuteDownloadTask(ctx context.Context, workerId uint, task *ddownload.DownloadTask) error
}
