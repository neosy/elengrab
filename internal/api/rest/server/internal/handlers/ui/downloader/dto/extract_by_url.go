package dto

type ExtractPlaylistByURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type ExtractPlaylistByURLResponse struct {
	Items []PlaylistItem `json:"items"`
}

type PlaylistItem struct {
	Index int `json:"index"`

	URL   string `json:"url"`
	Title string `json:"title"`

	DurationSec  int    `json:"durationSec"`
	DurationText string `json:"durationText"`
}
