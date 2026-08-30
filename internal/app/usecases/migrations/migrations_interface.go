package migrations

import "context"

type Migrations interface {
	// RunRequiredMigrations executes all required migrations.
	// The application cannot continue until these migrations complete successfully.
	RunRequiredMigrations(ctx context.Context) error

	// RunDeferredMigrations executes deferred migrations.
	// These migrations are not required for startup and can be run after the application is ready.
	RunDeferredMigrations(ctx context.Context) error
}
