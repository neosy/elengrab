package pworkers

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type DownloadTaskRunner interface {
	ExecuteDownloadTask(ctx context.Context, workerId uint, task *ddownload.DownloadTask) error
}

type MaintenanceRunner interface {
	UpdateHash(ctx context.Context) error
	DeleteDuplicates(ctx context.Context) error
	DeleteMissingFiles(ctx context.Context, enableMoveUnmatchedFiles bool) error
	DeleteFailedDownloads(ctx context.Context) error
}
