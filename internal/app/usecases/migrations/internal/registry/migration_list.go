package registry

import dtypes "github.com/neosy/elengrab/internal/domain/types"

type MigrationList struct {
	items []*MigrationID
}

func NewMigrationList() MigrationList {
	return MigrationList{}
}

func (m *MigrationList) Add(id string, migrationType dtypes.MigrationType, run MigrationRunner) {
	m.items = append(m.items, &MigrationID{
		id:            id,
		migrationType: migrationType,
		run:           run,
	})
}

func (m *MigrationList) Items() []*MigrationID {
	return m.items
}
