package extractinfo

import (
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

func processVKVideo(info *idto.ExtractInfo) {
	if info.ChannelID == "" && info.UploaderID != "" {
		info.ChannelID = strings.TrimPrefix(info.UploaderID, "-")
	}

	info.ParsedChannel.Title = info.Uploader
	info.ParsedChannel.Username = ""
	info.ParsedChannel.UsernameURL = ""
}
