package extractinfo

import (
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

func processX(info *idto.ExtractInfo) {
	if info.ChannelURL == "" && info.UploaderURL != "" {
		info.ChannelURL = info.UploaderURL
	}

	info.ParsedChannel.Title = info.Uploader
	info.ParsedChannel.Username = info.UploaderID
	info.ParsedChannel.UsernameURL = info.UploaderURL
}
