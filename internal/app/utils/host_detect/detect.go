package hostdetect

import dtypes "github.com/neosy/elengrab/internal/domain/types"

var matchers = []struct {
	platform   dtypes.MediaPlatform
	isPlatform isPlatform
}{
	{dtypes.MediaPlatformYouTube, YouTube},
	{dtypes.MediaPlatformFacebook, Facebook},
	{dtypes.MediaPlatformInstagram, Instagram},
	{dtypes.MediaPlatformTwitch, Twitch},
	{dtypes.MediaPlatformVimeo, Vimeo},
	{dtypes.MediaPlatformTikTok, TikTok},
	{dtypes.MediaPlatformRutube, Rutube},
}

func Detect(rawURL string) dtypes.MediaPlatform {
	for _, m := range matchers {
		if m.isPlatform(rawURL) {
			return m.platform
		}
	}
	return dtypes.MediaPlatformNone
}
