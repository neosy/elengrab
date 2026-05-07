package helper

import (
	"encoding/json"
	"net/url"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

type embedInfoResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorName  string `json:"author_name"`
	AuthorURL   string `json:"author_url"`
	Thumbnail   string `json:"thumbnail_url"`
	Provider    string `json:"provider_name"`
}

func FetchYoutubeInfoFast(youtubeUrl string) (*embedInfoResponse, error) {
	apiURL := "https://www.youtube.com/oembed?format=json&url=" + url.QueryEscape(youtubeUrl)

	client := httpx.NewClient(
		httpx.ClientOptionWithTimeout(consts.FetchTitleTimeout),
	)

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
