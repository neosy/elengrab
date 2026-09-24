package extractinfo

import (
	"fmt"
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

func processTwitch(info *idto.ExtractInfo) {
	usernameURL := fmt.Sprintf("https://www.twitch.tv/%s/", strings.ToLower(info.Channel))

	if info.ChannelURL == "" {
		info.ChannelURL = usernameURL
	}

	info.ParsedChannel.Title = info.Channel
	info.ParsedChannel.Username = info.Channel
	info.ParsedChannel.UsernameURL = usernameURL
}
