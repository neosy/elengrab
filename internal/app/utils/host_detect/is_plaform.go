package hostdetect

import (
	"net/url"
	"strings"
)

type isPlatform func(rawURL string) bool

var (
	YouTube = makeIsPlatform("youtube.com", map[string]struct{}{
		"youtube.com": {}, "www.youtube.com": {}, "m.youtube.com": {},
		"youtu.be": {}, "www.youtu.be": {}, "music.youtube.com": {},
	})

	Twitch = makeIsPlatform("twitch.tv", map[string]struct{}{
		"twitch.tv": {}, "www.twitch.tv": {}, "m.twitch.tv": {},
		"clips.twitch.tv": {}, "player.twitch.tv": {},
	})

	Vimeo = makeIsPlatform("vimeo.com", map[string]struct{}{
		"vimeo.com": {}, "www.vimeo.com": {}, "player.vimeo.com": {},
	})

	TikTok = makeIsPlatform("tiktok.com", map[string]struct{}{
		"tiktok.com": {}, "www.tiktok.com": {}, "m.tiktok.com": {},
		"vm.tiktok.com": {},
	})

	Facebook = makeIsPlatform("facebook.com", map[string]struct{}{
		"facebook.com": {}, "www.facebook.com": {},
		"m.facebook.com": {}, "fb.watch": {},
	})

	Instagram = makeIsPlatform("instagram.com", map[string]struct{}{
		"instagram.com": {}, "www.instagram.com": {}, "m.instagram.com": {},
		"instagr.am":      {},
		"l.instagram.com": {},
		"www.instagr.am":  {},
	})

	Rutube = makeIsPlatform("rutube.ru", map[string]struct{}{
		"rutube.ru": {}, "www.rutube.ru": {}, "m.rutube.ru": {},
		"rutube.su": {}, "www.rutube.su": {},
	})
)

func makeIsPlatform(rootDomain string, urls map[string]struct{}) isPlatform {
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
