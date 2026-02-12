package downloader

import "time"

const (
	logoUpdateInterval    = 24 * time.Hour
	channelUpdateInterval = 30 * 24 * time.Hour
	getHTMLTimeout        = 3 * time.Second
)
