package ytdlpsrv

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type embedInfoResponse struct {
	Title      string `json:"title"`
	AuthorName string `json:"author_name"`
	AuthorURL  string `json:"author_url"`
	Thumbnail  string `json:"thumbnail_url"`
	Provider   string `json:"provider_name"`
}

func (srv *YtDlpService) getInfoFast(youtubeUrl string) (*embedInfoResponse, error) {
	apiURL := "https://www.youtube.com/oembed?format=json&url=" + url.QueryEscape(youtubeUrl)

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data embedInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
