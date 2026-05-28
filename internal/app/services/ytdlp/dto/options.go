package dto

import "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"

type Option func(opt *Options)

// Options holds configuration options for YtDlpService
type Options struct {
	// Number of fragments of a dash/hlsnative video that should be downloaded concurrently (default is 1)
	ConcurrentFragments uint8

	// CookiesDir path for storing cookies
	CookiesDir string
	// AllowCookies allow cookies for YouTube, Instagram, etc
	AllowCookies bool
}

// NewOptions sets default values for Options fields
func NewOptions() Options {
	return Options{
		ConcurrentFragments: consts.ConcurrentFragmentsDefault,
	}
}

func WithConcurrentFragments(value uint8) Option {
	return func(opt *Options) {
		opt.ConcurrentFragments = value
	}
}

func WithCookiesDir(value string) Option {
	return func(opt *Options) {
		opt.CookiesDir = value
	}
}

func WithAllowCookies(value bool) Option {
	return func(opt *Options) {
		opt.AllowCookies = value
	}
}
