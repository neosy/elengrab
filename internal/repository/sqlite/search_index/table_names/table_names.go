package tablenames

const (
	MediaSourcesIndex = "media_sources_index"
)

var tableNames = []string{
	MediaSourcesIndex,
}

func TableNames() []string {
	return tableNames
}
