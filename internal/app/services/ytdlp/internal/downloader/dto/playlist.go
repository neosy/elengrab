package idto

type ExtractPlaylistItem struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	WebpageURL    string `json:"webpage_url"`
	PlaylistIndex int    `json:"playlist_index"`
}

type Playlist struct {
	Items []PlaylistItem
}

type PlaylistItem struct {
	URL   string
	Title string
}

func (p *Playlist) Count() int {
	if p == nil {
		return 0
	}

	return len(p.Items)
}
