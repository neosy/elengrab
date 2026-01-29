package downloader

import "time"

const (
	ytDlpTempDir  = ".yt-dlp"
	ytDlpCacheDir = ".yt-dlp/cache"

	channelAvatarTimeout = 5 * time.Second
	ytDlpTimeout         = 2 * time.Hour
	ytDlpRetryDelay      = 2 * time.Second
)
