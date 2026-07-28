package tablenames

const (
	MediaWatchEvents    = "media_watch_events"
	MediaWatchChunks    = "media_watch_chunks"
	MediaUserWatchStats = "media_user_watch_stats"
	MediaWatchStats     = "media_watch_stats"
	MediaWatchPositions = "media_watch_positions"
)

var tableNames = []string{
	MediaWatchEvents,
	MediaWatchChunks,
	MediaUserWatchStats,
	MediaWatchStats,
	MediaWatchPositions,
}

func TableNames() []string {
	return tableNames
}
