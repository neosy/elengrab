package extractinfo

import (
	"fmt"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

func processRuTube(info *idto.ExtractInfo) {
	if info.ChannelID == "" && info.UploaderID != "" {
		info.ChannelID = info.UploaderID
	}

	usernameURL := fmt.Sprintf("https://rutube.ru/channel/%s/", info.UploaderID)

	if info.ChannelURL == "" {
		info.ChannelURL = usernameURL
	}

	info.ParsedChannel.Title = info.Uploader
	info.ParsedChannel.Username = ""
	info.ParsedChannel.UsernameURL = usernameURL
}
