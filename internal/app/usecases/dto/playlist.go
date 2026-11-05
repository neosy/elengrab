package dto

type Playlist struct {
	Items []PlaylistItem
}

type PlaylistItem struct {
	URL   string
	Title string
}
