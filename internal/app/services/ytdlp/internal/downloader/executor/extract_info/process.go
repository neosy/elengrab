package extractinfo

import (
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func Process(mediaURL string, info *idto.ExtractInfo) {
	platformType := hostdetect.DetectPlatformType(mediaURL)

	switch platformType {
	case dtypes.MediaPlatformTypeInstagram:
		processInstagram(info)
	case dtypes.MediaPlatformTypeYouTube:
		processYouTube(info)
	case dtypes.MediaPlatformTypeTikTok:
		processTikTok(info)
	case dtypes.MediaPlatformTypeTwitch:
		processTwitch(info)
	case dtypes.MediaPlatformTypeX:
		processX(info)
	case dtypes.MediaPlatformTypeReddit:
		processReddit(mediaURL, info)
	case dtypes.MediaPlatformTypeRutube:
		processRuTube(info)
	case dtypes.MediaPlatformTypeVKVideo:
		processVKVideo(mediaURL, info)
	default:
		if info.ChannelURL == "" && info.UploaderURL != "" {
			info.ChannelURL = info.UploaderURL
		}

		info.ParsedChannel.Title = info.Uploader
	}
}
