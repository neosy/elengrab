package base

import (
	"context"
	"fmt"
)

func (m *Migrations) RunMigrations(ctx context.Context) error {
	var hasError = false

	for _, migration := range m.migrationList.Items() {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		default:
		}

		exists, err := m.Usecases().DownloadMigration.Exists(ctx, migration.ID())
		if err != nil {
			hasError = true
			continue
		}

		if exists {
			continue
		}

		m.logger.Info("Start data migration...", "id", migration.ID())

		done, err := migration.Run(ctx)
		if err != nil {
			m.logger.Warn("Failed data migration process '%s'", "id", migration.ID(), "error", err)
			hasError = true
			continue
		}

		if done {
			err := m.MarkMigration(ctx, migration.ID())
			if err != nil {
				hasError = true
			}
		}

		m.logger.Info("Data migration completed", "id", migration.ID())
	}

	if hasError {
		return fmt.Errorf("errors in the 'downloader usecase' migrations")
	}

	return nil
}
