package pworkers

import "context"

type MaintenanceRunner interface {
	BackupDatabase(ctx context.Context) error
}
