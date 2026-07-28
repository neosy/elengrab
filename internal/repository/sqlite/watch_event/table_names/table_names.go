package tablenames

const (
	MediaWatchEvents        = "media_watch_events"
	MediaUserWatchChunks    = "media_user_watch_chunks"
	MediaUserWatchStats     = "media_user_watch_stats"
	MediaWatchStats         = "media_watch_stats"
	MediaUserWatchPositions = "media_user_watch_positions"
)

var tableNames = []string{
	MediaWatchEvents,
	MediaUserWatchChunks,
	MediaUserWatchStats,
	MediaWatchStats,
	MediaUserWatchPositions,
}

func TableNames() []string {
	return tableNames
}
