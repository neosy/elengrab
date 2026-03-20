package pworkers

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type DownloadTaskRunner interface {
	ExecuteDownloadTask(ctx context.Context, workerId uint64, task *ddownload.DownloadTask) error
	UpdateSystemInfo()
}
