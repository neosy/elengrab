package pworkers

import "context"

type DownloadMaintenanceRunner interface {
	UpdateHash(ctx context.Context) error
	DeleteDuplicates(ctx context.Context) error
	DeleteMissingFiles(ctx context.Context, enableMoveUnmatchedFiles bool) error
	DeleteFailedDownloads(ctx context.Context) error
}
