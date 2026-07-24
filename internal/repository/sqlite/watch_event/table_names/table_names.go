package tablenames

const (
	MediaWatchEvents = "media_watch_events"
	MediaWatchChunks = "media_watch_chunks"
	MediaWatchStats  = "media_watch_stats"
)

var tableNames = []string{
	MediaWatchEvents,
	MediaWatchChunks,
	MediaWatchStats,
}

func TableNames() []string {
	return tableNames
}
