package tablenames

const (
	MediaWatchEvents   = "media_watch_events"
	MediaWatchChunks   = "media_watch_chunks"
	MediaWatchStats    = "media_watch_stats"
	MediaWatchPositions = "media_watch_positions"
)

var tableNames = []string{
	MediaWatchEvents,
	MediaWatchChunks,
	MediaWatchStats,
	MediaWatchPositions,
}

func TableNames() []string {
	return tableNames
}
