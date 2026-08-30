package maintenance

import "context"

type Maintenance interface {
	StartupDatabase(ctx context.Context) error

	BackupDatabase(ctx context.Context) error
	FlushWAL() error
}
