package dlmigration

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *DownloadMigration) MarkMigration(
	ctx context.Context,
	migrationID string,
	migrationType dtypes.MigrationType,
) error {
	migration := &ddownload.DataMigration{
		MigrationID:   migrationID,
		MigrationType: migrationType,
	}

	err := m.Insert(ctx, migration)
	if err != nil {
		return err
	}

	return nil
}
