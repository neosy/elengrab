package pworkers

import "context"

type MaintenanceRunner interface {
	BackupDatabase(ctx context.Context) error
	FlushWAL() error
	StartupDatabase(ctx context.Context) error
}
