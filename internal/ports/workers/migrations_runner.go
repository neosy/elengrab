package pworkers

import "context"

type MigrationsRunner interface {
	ExecuteMigrations(ctx context.Context) error
}
