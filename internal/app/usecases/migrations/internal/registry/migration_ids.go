package registry

import (
	"context"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type (
	MigrationRunner func(context.Context) error

	MigrationID struct {
		id            string
		migrationType dtypes.MigrationType
		run           MigrationRunner
	}
)

func (i *MigrationID) ID() string {
	return i.id
}

func (i *MigrationID) MigrationType() dtypes.MigrationType {
	return i.migrationType
}

func (i *MigrationID) Run(ctx context.Context) error {
	return i.run(ctx)
}
