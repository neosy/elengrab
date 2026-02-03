package logofetcher

import "time"

const (
	limitImage           = 512 << 10 // 512 KB
	limitHTML            = 512 << 10 // 512 KB
	downloadImageTimeout = 3 * time.Second
	getHTMLTimeout       = 3 * time.Second
)
