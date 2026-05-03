package hostdetect

import dtypes "github.com/neosy/elengrab/internal/domain/types"

var matchers = []struct {
	host   dtypes.MediaHost
	isHost isHost
}{
	{dtypes.MediaHostYouTube, YouTube},
	{dtypes.MediaHostFacebook, Facebook},
	{dtypes.MediaHostInstagram, Instagram},
	{dtypes.MediaHostTwitch, Twitch},
	{dtypes.MediaHostVimeo, Vimeo},
	{dtypes.MediaHostTikTok, TikTok},
	{dtypes.MediaHostRutube, Rutube},
}

func Detect(rawURL string) dtypes.MediaHost {
	for _, m := range matchers {
		if m.isHost(rawURL) {
			return m.host
		}
	}
	return dtypes.MediaHostNone
}
