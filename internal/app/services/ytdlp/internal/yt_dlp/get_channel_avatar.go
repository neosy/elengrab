package ytdlp

import (
	"encoding/json"
	"fmt"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/yt_dlp/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/yt_dlp/helper"
	"github.com/neosy/elengrab/pkg/nfasthttp"
)

// getChannelAvatar fetches the HTML of a YouTube channel page,
// extracts the avatar JSON block, and returns all avatar URLs.
// url: full URL of the YouTube channel
func (y *YTDlp) getChannelAvatar(url string) ([]idto.AvatarSource, error) {
	body, err := nfasthttp.GetHTML(
		url,
		nfasthttp.ClientOptionWithreadBufferSize(64*1024),
		nfasthttp.ClientOptionWithTimeout(channelAvatarTimeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get html: %w", err)
	}

	html := string(body)

	// Full key path to the avatar sources inside decoratedAvatarViewModel
	fullKey := `"decoratedAvatarViewModel":{"avatar":{"avatarViewModel":{"image":{"sources":[`
	jsonArray, err := helper.ExtractJSONArray(html, fullKey)
	if err != nil {
		return nil, err
	}

	var sources []idto.AvatarSource
	if err := json.Unmarshal([]byte(jsonArray), &sources); err != nil {
		return nil, err
	}

	for i, src := range sources {
		raw, format, err := nfasthttp.GetImage(
			src.URL,
			nfasthttp.ClientOptionWithreadBufferSize(64*1024),
			nfasthttp.ClientOptionWithTimeout(channelAvatarTimeout),
		)
		if err != nil {
			continue
		}
		sources[i].Raw = raw
		sources[i].Format = format
	}

	return sources, nil
}
