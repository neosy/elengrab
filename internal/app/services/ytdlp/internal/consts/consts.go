package consts

import (
	"slices"
	"time"
)

const (
	YtDlpName = "yt-dlp"
	DenoName  = "deno"

	YtDlpTempDir        = ".yt-dlp"
	YtDlpCacheDir       = ".yt-dlp/cache"
	YtDlpFormatCacheDir = ".yt-dlp/format-cache"

	ChannelAvatarTimeout = 5 * time.Second
	FetchTitleTimeout    = 3 * time.Second
	YtDlpTimeout         = 2 * time.Hour
	YtDlpRetryDelay      = 2 * time.Second

	FetchImageLimit   = 4096 << 10 // 4096 KB
	FetchImageTimeout = 10 * time.Second

	MediaInfoExtractionTimeout = 5 * time.Second
	ThumbnailExtractionTimeout = 8 * time.Second

	FormatCacheTTL = 2 * time.Hour

	ConcurrentFragmentsDefault = 5
	MaxTitleLengthInFilename   = 100
)

var (
	shortYoutubeThumbnailURLTemplates = [...]string{
		"https://i.ytimg.com/vi/%s/oardefault.jpg",
		"https://i.ytimg.com/vi/%s/oar2.jpg",
	}
)

func ShortYoutubeThumbnailURLTemplates() []string {
	return slices.Clone(shortYoutubeThumbnailURLTemplates[:])
}
