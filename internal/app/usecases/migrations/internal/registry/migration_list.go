package registry

type MigrationList struct {
	items []*MigrationID
}

func NewMigrationList() MigrationList {
	return MigrationList{}
}

func (m *MigrationList) Add(id string, run MigrationRunner) {
	m.items = append(m.items, &MigrationID{
		id:  id,
		run: run,
	})
}

func (m *MigrationList) Items() []*MigrationID {
	return m.items
}
