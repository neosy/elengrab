package migrations

func (m *migrations) initMigrations() {
	m.requiredMigrationIDs.addMigration("01_move_downloads_to_storage", m.moveDownloadsToStorage)

	m.deferredMigrationIDs.addMigration("02_fill_media_info", m.fillMediaInfo)
	m.deferredMigrationIDs.addMigration("03_fill_thumbnails", m.fillThumbnails)
	m.deferredMigrationIDs.addMigration("04_fill_title_for_instagram", m.fillTitleForInstagram)
	m.deferredMigrationIDs.addMigration("05_fill_media_description", m.fillMediaDescription)
}
