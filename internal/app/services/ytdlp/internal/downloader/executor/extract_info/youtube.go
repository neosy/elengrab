package extractinfo

import (
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

func processYouTube(info *idto.ExtractInfo) {
	if info.ChannelURL == "" && info.UploaderURL != "" {
		info.ChannelURL = info.UploaderURL
	}

	info.ParsedChannel.Title = info.Uploader
	info.ParsedChannel.Username = strings.TrimPrefix(info.UploaderID, "@")
	info.ParsedChannel.UsernameURL = info.UploaderURL
}
