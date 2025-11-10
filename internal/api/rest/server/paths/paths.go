package httppaths

// Index
const (
	PathIndex = "/"
)

// UI
const (
	// Groups
	GroupStatic     = "/ui/static"
	GroupDownloader = "/ui/downloader"

	// Paths
	PathGrab     = "/grab"
	PathDownload = "/download"
	PathFiles    = "/{filepath:*}"
)
