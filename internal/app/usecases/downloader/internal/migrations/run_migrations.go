package migrations

import (
	"context"
	"fmt"
)

func (m *migrations) RunRequiredMigrations(ctx context.Context) error {
	return m.runMigrations(ctx, m.requiredMigrationList)
}

func (m *migrations) RunDeferredMigrations(ctx context.Context) error {
	return m.runMigrations(ctx, m.deferredMigrationList)
}

func (m *migrations) runMigrations(ctx context.Context, migrationList migrationList) error {
	var hasError = false

	for _, migration := range migrationList.items {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		default:
		}

		exists, err := m.usecases.downloadMigration.Exists(ctx, migration.id)
		if err != nil {
			hasError = true
			continue
		}

		if exists {
			continue
		}

		m.logger.Info("Start data migration...", "id", migration.id)

		done, err := migration.run(ctx)
		if err != nil {
			m.logger.Warn("Failed data migration process '%s'", "id", migration.id, "error", err)
			hasError = true
			continue
		}

		if done {
			err := m.markMigration(ctx, migration.id)
			if err != nil {
				hasError = true
			}
		}

		m.logger.Info("Data migration completed", "id", migration.id)
	}

	if hasError {
		return fmt.Errorf("errors in the 'downloader usecase' migrations")
	}

	return nil
}
