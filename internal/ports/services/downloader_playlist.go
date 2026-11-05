package pservices

// DownloaderPlaylist represents a collection of playlist items, each containing its URL and title.
type DownloaderPlaylist struct {
	Items []DownloaderPlaylistItem
}

// DownloaderPlaylistItem represents a single item in a playlist, containing its URL and title.
type DownloaderPlaylistItem struct {
	URL   string
	Title string
}

func (p *DownloaderPlaylist) Count() int {
	if p == nil {
		return 0
	}

	return len(p.Items)
}
