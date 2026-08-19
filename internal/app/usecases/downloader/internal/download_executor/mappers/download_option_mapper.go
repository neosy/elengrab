package mappers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

func (m *Mappers) MapDownloadOptionsDomainToService(options *ddownload.DownloadOptions) *dservices.DownloadOptions {
	opt := dservices.NewDefaultDownloadOptions()

	opt.FormatType = options.FormatType
	opt.VideoFormat = options.VideoFormat
	opt.VideoCodec = options.VideoCodec
	opt.VideoResolution = options.VideoResolution
	opt.AudioFormat = options.AudioFormat
	opt.Filename = options.Filename
	opt.VideoQuality = options.VideoQuality
	opt.AudioQuality = options.AudioQuality

	return opt
}
