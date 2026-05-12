package migrations

import (
	"context"
	"fmt"
)

func (m *migrations) ExecuteMigrations(ctx context.Context) error {
	var hasError = false

	for id, migration := range m.migrationIDs {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		default:
		}

		exists, err := m.usecases.downloadMigration.Exists(ctx, id)
		if err != nil {
			hasError = true
			continue
		}

		if exists {
			continue
		}

		m.logger.Info("Start data migration...", "id", id)

		done, err := migration.run(ctx)
		if err != nil {
			m.logger.Warn("Failed data migration process '%s'", "id", id, "error", err)
			hasError = true
			continue
		}

		if done {
			err := m.markMigration(ctx, id)
			if err != nil {
				hasError = true
			}
		}

		m.logger.Info("Data migration completed", "id", id)
	}

	if hasError {
		return fmt.Errorf("errors in the 'downloader usecase' migrations")
	}

	return nil
}
