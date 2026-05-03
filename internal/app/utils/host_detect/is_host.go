package hostdetect

import (
	"net/url"
	"strings"
)

type isHost func(rawURL string) bool

var (
	YouTube = makeIsHost("youtube.com", map[string]struct{}{
		"youtube.com": {}, "www.youtube.com": {}, "m.youtube.com": {},
		"youtu.be": {}, "www.youtu.be": {}, "music.youtube.com": {},
	})

	Twitch = makeIsHost("twitch.tv", map[string]struct{}{
		"twitch.tv": {}, "www.twitch.tv": {}, "m.twitch.tv": {},
		"clips.twitch.tv": {}, "player.twitch.tv": {},
	})

	Vimeo = makeIsHost("vimeo.com", map[string]struct{}{
		"vimeo.com": {}, "www.vimeo.com": {}, "player.vimeo.com": {},
	})

	TikTok = makeIsHost("tiktok.com", map[string]struct{}{
		"tiktok.com": {}, "www.tiktok.com": {}, "m.tiktok.com": {},
		"vm.tiktok.com": {},
	})

	Facebook = makeIsHost("facebook.com", map[string]struct{}{
		"facebook.com": {}, "www.facebook.com": {},
		"m.facebook.com": {}, "fb.watch": {},
	})

	Instagram = makeIsHost("instagram.com", map[string]struct{}{
		"instagram.com": {}, "www.instagram.com": {}, "m.instagram.com": {},
		"instagr.am":      {},
		"l.instagram.com": {},
		"www.instagr.am":  {},
	})

	Rutube = makeIsHost("rutube.ru", map[string]struct{}{
		"rutube.ru": {}, "www.rutube.ru": {}, "m.rutube.ru": {},
		"rutube.su": {}, "www.rutube.su": {},
	})
)

func makeIsHost(rootDomain string, urls map[string]struct{}) isHost {
	return func(rawURL string) bool {
		u, err := url.Parse(rawURL)
		if err != nil {
			return false
		}

		host := strings.ToLower(u.Hostname())

		if _, ok := urls[host]; ok {
			return true
		}

		// subdomains like foo.youtube.com
		return strings.HasSuffix(host, "."+rootDomain)
	}
}
