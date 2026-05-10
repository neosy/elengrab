package migrations

import "context"

type (
	migrationRunner func(context.Context) (bool, error)
	migrationID     struct {
		id  string
		run migrationRunner
	}
	migrationIDMap map[string]*migrationID
)

func (m migrationIDMap) addMigration(id string, run migrationRunner) {
	m[id] = &migrationID{
		id:  id,
		run: run,
	}
}
