package mappers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
)

func (m *Mappers) MapDownloadOptionsDomainToService(options *ddownload.DownloadOptions) *dservices.DownloadOptions {
	return &dservices.DownloadOptions{
		FormatType:             options.FormatType,
		VideoFormat:            options.VideoFormat,
		VideoCodec:             options.VideoCodec,
		VideoResolution:        options.VideoResolution,
		AudioFormat:            options.AudioFormat,
		Filename:               options.Filename,
		IncludeTitleInFilename: true,
		VideoQuality:           options.VideoQuality,
		AudioQuality:           options.AudioQuality,
	}
}
