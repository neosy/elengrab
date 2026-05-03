package dlmigration

import "context"

func (uc *DataMigration) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.dataMigrationRep.Tx(ctx, fn)
}
