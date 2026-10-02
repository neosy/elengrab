package extractinfo

import (
	"net/url"
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

func processVKVideo(mediaURL string, info *idto.ExtractInfo) {
	if info.ChannelID == "" && info.UploaderID != "" {
		info.ChannelID = strings.TrimPrefix(info.UploaderID, "-")
	}

	if info.ChannelID != "" {
		username := "club" + info.ChannelID
		info.ParsedChannel.Username = username
	}

	if info.ParsedChannel.Username != "" {
		parsedURL, err := url.Parse(mediaURL)
		if err == nil {
			info.ParsedChannel.UsernameURL = "https://" + parsedURL.Host + "/@" + info.ParsedChannel.Username
		}
	}

	if info.ChannelURL == "" {
		info.ChannelURL = info.ParsedChannel.UsernameURL
	}

	info.ParsedChannel.Title = info.Uploader
}
