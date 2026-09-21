package required

func (m *migrations) initMigrations() {
	m.Add("01_move_downloads_to_storage", m.moveDownloadsToStorage)
	m.Add("02_change_watch_chunk_size", m.changeWatchChunkSizeOnce)
	m.Add("03_rebuild_watch_stats", m.rebuildWatchStatsOnce)
	m.Add("04_build_search_index", m.buildSearchIndexOnce)
	m.Add("05_add_multi_platform_channels", m.addMultiPlatformChannels)
	m.Add("06_fill_search_index_channel_id", m.fillSearchIndexChannelID)
}
