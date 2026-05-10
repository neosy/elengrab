package migrations

func (m *migrations) initMigrations() {
	m.migrationIDs.addMigration("01_move_downloads_to_storage", m.moveDownloadsToStorage)
	m.migrationIDs.addMigration("02_fetch_thumbnails", m.fetchThumbnails)
}
