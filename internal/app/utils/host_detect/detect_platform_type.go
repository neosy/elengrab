package hostdetect

import dtypes "github.com/neosy/elengrab/internal/domain/types"

var matchers = []struct {
	platform   dtypes.MediaPlatformType
	isPlatform isPlatform
}{
	{dtypes.MediaPlatformTypeYouTube, YouTube},
	{dtypes.MediaPlatformTypeFacebook, Facebook},
	{dtypes.MediaPlatformTypeInstagram, Instagram},
	{dtypes.MediaPlatformTypeTwitch, Twitch},
	{dtypes.MediaPlatformTypeVimeo, Vimeo},
	{dtypes.MediaPlatformTypeTikTok, TikTok},
	{dtypes.MediaPlatformTypeRutube, Rutube},
}

func DetectPlatformType(rawURL string) dtypes.MediaPlatformType {
	for _, m := range matchers {
		if m.isPlatform(rawURL) {
			return m.platform
		}
	}
	return dtypes.MediaPlatformTypeNone
}
