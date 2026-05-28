package idto

import dtypes "github.com/neosy/elengrab/internal/domain/types"

type DLOptions struct {
	FileName               string
	IncludeTitleInFilename bool

	FormatType          dtypes.FormatType
	VideoFormat         dtypes.VideoFormat
	VideoCodec          dtypes.VideoCodec
	VideoResolution     dtypes.VideoResolution
	AudioFormat         dtypes.AudioFormat
	ConcurrentFragments uint8

	NeedsCookies   bool
	CookieFilePath string
}

func (opts *DLOptions) AllowCookies() bool {
	return opts.CookieFilePath != ""
}

func (opts *DLOptions) CookieFilePathIfNeeded() string {
	if opts.NeedsCookies == false {
		return ""
	}
	return opts.CookieFilePath
}
