package consts

import "time"

const (
	YtDlpName = "yt-dlp"
	DenoName  = "deno"

	YtDlpTempDir               = ".yt-dlp"
	YtDlpCacheDir              = ".yt-dlp/cache"
	YtDlpFormatCacheDir        = ".yt-dlp/format-cache"
	YtDlpYouTubeCookieFileName = "youtube-cookies.txt"

	ChannelAvatarTimeout = 5 * time.Second
	FetchTitleTimeout    = 3 * time.Second
	YtDlpTimeout         = 2 * time.Hour
	YtDlpRetryDelay      = 2 * time.Second

	FormatCacheTTL = 2 * time.Hour

	ConcurrentFragmentsDefault = 5
	MaxTitleLengthInFilename   = 100
)
