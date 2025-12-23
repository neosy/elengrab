package maintenance

func (m *Maintenance) FlushWAL() error {
	return m.database.FlushWAL()
}
