package pworkers

import "context"

type MigrationsRunner interface {
	RunDeferredMigrations(ctx context.Context) error
}
