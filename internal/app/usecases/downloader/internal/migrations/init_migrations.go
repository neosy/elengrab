package migrations

func (m *migrations) initMigrations() {
	m.migrationIDs.addMigration("01_move_downloads_to_storage", m.moveDownloadsToStorage)
	m.migrationIDs.addMigration("02_fill_media_info", m.fillMediaInfo)
	m.migrationIDs.addMigration("03_fill_thumbnails", m.fillThumbnails)
}
