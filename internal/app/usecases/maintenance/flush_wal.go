package maintenance

func (m *maintenance) FlushWAL() error {
	return m.repositories.FlushWAL()
}
