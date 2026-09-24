package extractinfo

import (
	"fmt"
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

func processInstagram(info *idto.ExtractInfo) {
	if info.Description != "" {
		first, _, _ := strings.Cut(info.Description, "\n")
		if first != "" {
			info.Title = first
		}
	}

	if info.ChannelID == "" && info.UploaderID != "" {
		info.ChannelID = info.UploaderID
	}

	usernameURL := fmt.Sprintf("https://www.instagram.com/%s/", info.Channel)

	if info.ChannelURL == "" {
		info.ChannelURL = usernameURL
	}

	info.ParsedChannel.Title = info.Uploader
	info.ParsedChannel.Username = info.Channel
	info.ParsedChannel.UsernameURL = usernameURL
}
