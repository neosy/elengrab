package tablenames

const (
	YoutubeChannels = "youtube_channels"
	SiteLogos       = "site_logos"
	Thumbnails      = "media_thumbnails"
)

var tableNames = []string{
	YoutubeChannels,
	SiteLogos,
	Thumbnails,
}

func TableNames() []string {
	return tableNames
}
