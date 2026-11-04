package pservices

import "time"

// DownloaderPlaylist represents a collection of playlist items, each containing its URL and title.
type DownloaderPlaylist struct {
	Items []DownloaderPlaylistItem
}

// DownloaderPlaylistItem represents a single item in a playlist, containing its URL and title.
type DownloaderPlaylistItem struct {
	ID    string
	Index int

	URL   string
	Title string

	Duration time.Duration
}

func (p *DownloaderPlaylist) Count() int {
	if p == nil {
		return 0
	}

	return len(p.Items)
}
