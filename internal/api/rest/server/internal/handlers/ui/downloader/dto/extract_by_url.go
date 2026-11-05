package dto

type ExtractPlaylistByURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type ExtractPlaylistByURLResponse struct {
	Items []PlaylistItem `json:"items"`
}

type PlaylistItem struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}
