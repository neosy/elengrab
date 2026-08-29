package dlmigration

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (m *DownloadMigration) MarkMigration(ctx context.Context, migrationID string) error {
	migration := &ddownload.DataMigration{
		MigrationID: migrationID,
	}

	err := m.Insert(ctx, migration)
	if err != nil {
		return err
	}

	return nil
}
