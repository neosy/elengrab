package pservices

// DownloaderMeta is the raw metadata reported by yt-dlp:
// title, description and the list of available formats.
type DownloaderMeta struct {
	Title       string
	Description string
	Formats     []DownloaderMetaMediaFormat
}
