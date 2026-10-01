package extractinfo

import (
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

func processTikTok(info *idto.ExtractInfo) {
	if info.ChannelID == "" && info.UploaderID != "" {
		info.ChannelID = info.UploaderID
	}

	if info.ChannelURL == "" && info.UploaderURL != "" {
		info.ChannelURL = info.UploaderURL
	}

	info.ParsedChannel.Title = info.Channel
	info.ParsedChannel.Username = info.Uploader
	info.ParsedChannel.UsernameURL = info.UploaderURL
}
