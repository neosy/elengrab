package mappers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func (m *Mappers) MapDownloadOptionsDomainToService(options *ddownload.DownloadOptions) *pservices.DownloadOptions {
	opt := pservices.NewDefaultDownloadOptions()

	opt.FormatType = options.FormatType
	opt.VideoFormat = options.VideoFormat
	opt.VideoCodec = options.VideoCodec
	opt.VideoResolution = options.VideoResolution
	opt.AudioFormat = options.AudioFormat
	opt.UniqueFileName = options.UniqueFileName
	opt.VideoQuality = options.VideoQuality
	opt.AudioQuality = options.AudioQuality

	return opt
}
