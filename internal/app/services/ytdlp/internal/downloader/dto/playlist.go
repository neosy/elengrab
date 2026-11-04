package idto

import "time"

type ExtractPlaylistItem struct {
	ID            string `json:"id"`
	PlaylistIndex int    `json:"playlist_index"`

	Title      string `json:"title"`
	URL        string `json:"url"`
	WebpageURL string `json:"webpage_url"`

	DurationSec  int    `json:"duration"`
	DurationText string `json:"duration_string"`
}

type Playlist struct {
	Items []PlaylistItem
}

type PlaylistItem struct {
	ID    string
	Index int

	URL   string
	Title string

	Duration time.Duration
}

func (p *Playlist) Count() int {
	if p == nil {
		return 0
	}

	return len(p.Items)
}
