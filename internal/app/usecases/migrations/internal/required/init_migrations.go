package required

func (m *migrations) initMigrations() {
	m.Add("01_move_downloads_to_storage", m.moveDownloadsToStorage)
	m.Add("02_change_watch_chunk_size", m.changeWatchChunkSizeOnce)
	m.Add("03_rebuild_watch_stats", m.rebuildWatchStatsOnce)
}
