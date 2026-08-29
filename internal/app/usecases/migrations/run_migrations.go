package migrations

import (
	"context"
)

// RunRequiredMigrations executes all required migrations.
// The application cannot continue until these migrations complete successfully.
func (m *migrations) RunRequiredMigrations(ctx context.Context) error {
	return m.required.RunMigrations(ctx)
}

// RunDeferredMigrations executes deferred migrations.
// These migrations are not required for startup and can be run after the application is ready.
func (m *migrations) RunDeferredMigrations(ctx context.Context) error {
	return m.deferred.RunMigrations(ctx)
}
