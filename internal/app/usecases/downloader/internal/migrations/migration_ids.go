package migrations

import "context"

type (
	migrationRunner func(context.Context) (bool, error)
	migrationID     struct {
		id  string
		run migrationRunner
	}
	migrationList struct {
		items []*migrationID
	}
)

func NewMigrationList() migrationList {
	return migrationList{}
}

func (m *migrationList) add(id string, run migrationRunner) {
	m.items = append(m.items, &migrationID{
		id:  id,
		run: run,
	})
}
