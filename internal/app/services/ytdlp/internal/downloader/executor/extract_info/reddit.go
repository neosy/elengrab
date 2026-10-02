package extractinfo

import (
	"net/url"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

func processReddit(mediaURL string, info *idto.ExtractInfo) {
	if info.ChannelURL == "" {
		info.ChannelURL = info.UploaderURL
	}

	if info.Uploader != "" {
		parsedURL, err := url.Parse(mediaURL)
		if err == nil {
			info.ParsedChannel.UsernameURL = "https://" + parsedURL.Host + "/user/" + info.Uploader
		}
	}

	if info.ChannelURL == "" {
		info.ChannelURL = info.ParsedChannel.UsernameURL
	}

	info.ChannelID = info.Uploader
	info.ParsedChannel.Title = info.Uploader
	info.ParsedChannel.Username = info.Uploader
}
