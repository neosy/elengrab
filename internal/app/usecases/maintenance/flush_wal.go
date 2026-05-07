package maintenance

func (m *Maintenance) FlushWAL() error {
	return m.repositories.FlushWAL()
}
