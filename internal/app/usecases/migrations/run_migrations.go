package migrations

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/migrations/internal/deferred"
	"github.com/neosy/elengrab/internal/app/usecases/migrations/internal/required"
)

// RunRequiredMigrations executes all required migrations.
// The application cannot continue until these migrations complete successfully.
func (m *migrations) RunRequiredMigrations(ctx context.Context) error {
	migrations := required.NewMigrations(m.logger, m.downloadsStorage, m.usecases, m.services)
	return migrations.RunMigrations(ctx)
}

// RunDeferredMigrations executes deferred migrations.
// These migrations are not required for startup and can be run after the application is ready.
func (m *migrations) RunDeferredMigrations(ctx context.Context) error {
	migrations := deferred.NewMigrations(m.logger, m.downloadsStorage, m.usecases, m.services)
	return migrations.RunMigrations(ctx)
}
