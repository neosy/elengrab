package pworkers

import "context"

// DownloadMaintenanceRunner defines the interface for a download maintenance runner.
type DownloadMaintenanceRunner interface {
	UpdateHash(ctx context.Context) error
	DeleteDuplicates(ctx context.Context) error
	DeleteMissingDownloads(ctx context.Context, enableMoveUnmatchedDownloads bool) error
	DeleteFailedDownloads(ctx context.Context) error
}
