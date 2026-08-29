package dlmigration

import "context"

func (uc *DownloadMigration) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.dataMigrationRepo().Tx(ctx, fn)
}
