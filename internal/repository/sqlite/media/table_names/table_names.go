package tablenames

const (
	Channels   = "channels"
	SiteLogos  = "site_logos"
	Thumbnails = "media_thumbnails"
)

var tableNames = []string{
	Channels,
	SiteLogos,
	Thumbnails,
}

func TableNames() []string {
	return tableNames
}
