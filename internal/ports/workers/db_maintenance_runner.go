package pworkers

import "context"

type DBMaintenanceRunner interface {
	BackupDatabase(ctx context.Context) error
	FlushWAL() error
	StartupDatabase(ctx context.Context) error
}
