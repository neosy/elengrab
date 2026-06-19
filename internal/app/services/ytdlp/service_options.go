package ytdlpsrv

import (
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

type ServiceOption func(opt *ServiceOptions)

// ServiceOptions holds configuration options for YtDlpService
type ServiceOptions struct {
	// Number of fragments of a dash/hlsnative video that should be downloaded concurrently (default is 1)
	ConcurrentFragments uint8

	// CookiesDir path for storing cookies
	CookiesDir string
	// AllowCookies allow cookies for YouTube, Instagram, etc
	AllowCookies bool
}

func (o ServiceOptions) toInternalOptions() idto.ServiceOptions {
	return idto.ServiceOptions{
		ConcurrentFragments: o.ConcurrentFragments,
		CookiesDir:          o.CookiesDir,
		AllowCookies:        o.AllowCookies,
	}
}

// NewServiceOptions sets default values for Options fields
func NewServiceOptions() ServiceOptions {
	return ServiceOptions{
		ConcurrentFragments: consts.ConcurrentFragmentsDefault,
	}
}

func WithConcurrentFragments(value uint8) ServiceOption {
	return func(opt *ServiceOptions) {
		opt.ConcurrentFragments = value
	}
}

func WithCookiesDir(value string) ServiceOption {
	return func(opt *ServiceOptions) {
		opt.CookiesDir = value
	}
}

func WithCookies(value bool) ServiceOption {
	return func(opt *ServiceOptions) {
		opt.AllowCookies = value
	}
}
