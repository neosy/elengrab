package base

import (
	"context"
	"fmt"
)

func (m *Migrations) RunMigrations(ctx context.Context) error {
	for _, migration := range m.migrationList.Items() {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		default:
		}

		exists, err := m.Usecases().DownloadMigration.Exists(ctx, migration.ID())
		if err != nil {
			return fmt.Errorf("check migration %s: %w", migration.ID(), err)
		}

		if exists {
			continue
		}

		m.logger.Info("Start data migration...", "id", migration.ID())

		err = migration.Run(ctx)
		if err != nil {
			m.logger.Warn("Failed data migration process", "id", migration.ID(), "error", err)
			return fmt.Errorf("run migration %s: %w", migration.ID(), err)
		}

		err = m.MarkMigration(ctx, migration.ID(), migration.MigrationType())
		if err != nil {
			return fmt.Errorf("mark migration %s: %w", migration.ID(), err)
		}

		m.logger.Info("Data migration completed", "id", migration.ID())
	}

	return nil
}
