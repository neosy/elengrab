package dto

import "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"

// Options holds configuration options for YtDlpService
type Options struct {
	// Number of fragments of a dash/hlsnative video that should be downloaded concurrently (default is 1)
	ConcurrentFragments uint8

	// CookiesDir path for storing cookies
	CookiesDir string
	// YoutubeAllowCookies allow cookies for YouTube
	YoutubeAllowCookies bool
}

// SetDefaults sets default values for Options fields if they are not set
// or if force is true
func (o *Options) SetDefaults(force bool) {
	if o.ConcurrentFragments == 0 || force {
		o.ConcurrentFragments = consts.ConcurrentFragmentsDefault
	}
}
