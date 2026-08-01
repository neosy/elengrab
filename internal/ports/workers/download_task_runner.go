package pworkers

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type DownloadTaskRunner interface {
	ExecuteDownloadTask(ctx context.Context, workerID uint64, task *ddownload.DownloadTask) error
	ExecuteRefreshMetadataTask(ctx context.Context, workerID uint64, task *ddownload.RefreshMetadataTask) error
	ExecuteCreateMediaWatchEvent(ctx context.Context, workerID uint64, req *dto.CreateMediaWatchEventRequest) error
	UpdateSystemInfo()
}
