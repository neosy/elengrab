package dto

import (
	"strings"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
)

type ExtractPlaylistByURLRequest struct {
	URL string `json:"url" validate:"required,url,max=2048"`
}

type ExtractPlaylistByURLResponse struct {
	Items []PlaylistItem `json:"items"`

	Count     int  `json:"count"`
	Truncated bool `json:"truncated"`
}

type PlaylistItem struct {
	Index int `json:"index"`

	URL   string `json:"url"`
	Title string `json:"title"`

	DurationSec  int    `json:"durationSec"`
	DurationText string `json:"durationText"`
}

func (req ExtractPlaylistByURLRequest) Validate() error {
	url := strings.TrimSpace(req.URL)
	if url == "" {
		return apierrors.ErrURLIsRequired
	}

	return nil
}
