package migrations

func (m *migrations) initMigrations() {
	m.requiredMigrationList.add("01_move_downloads_to_storage", m.moveDownloadsToStorage)
	m.requiredMigrationList.add("02_change_watch_chunk_size", m.changeWatchChunkSizeOnce)
	m.requiredMigrationList.add("03_rebuild_watch_stats", m.rebuildWatchStatsOnce)

	m.deferredMigrationList.add("02_fill_media_info", m.fillMediaInfo)
	m.deferredMigrationList.add("03_fill_thumbnails", m.fillThumbnails)
	m.deferredMigrationList.add("04_fill_title_for_instagram", m.fillTitleForInstagram)
	m.deferredMigrationList.add("05_fill_media_description", m.fillMediaDescription)
	m.deferredMigrationList.add("06_fill_media_duration", m.fillMediaDuration)
}
