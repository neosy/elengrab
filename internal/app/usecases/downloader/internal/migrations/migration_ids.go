package migrations

import "context"

type (
	migrationRunner func(context.Context) (bool, error)
	migrationID     struct {
		id  string
		run migrationRunner
	}
	migrationIDs []*migrationID
)

func (m migrationIDs) addMigration(id string, run migrationRunner) {
	m = append(m, &migrationID{
		id:  id,
		run: run,
	})
}
