package registry

import "context"

type (
	MigrationRunner func(context.Context) (bool, error)

	MigrationID struct {
		id  string
		run MigrationRunner
	}
)

func (i *MigrationID) ID() string {
	return i.id
}

func (i *MigrationID) Run(ctx context.Context) (bool, error) {
	return i.run(ctx)
}
