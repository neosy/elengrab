package helper

import (
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfile "github.com/neosy/elengrab/internal/pkg/file"
)

// PrepareDownloadOptions prepare download options with defaults and user overrides
// return dlOptions, downloadDir, fileName, includeTitleInFilename
func PrepareDownloadOptions(
	dlDir string,
	concurrentFragments uint8,
	options *dservices.DownloadOptions,
) (idto.DLOptions, string, string, bool) {
	// Set default values
	dlOptions := idto.DLOptions{
		FormatType:          formatTypeDefault,
		VideoFormat:         videoFormatDefault,
		VideoCodec:          videoCodecDefault,
		VideoResolution:     videoResolutionDefault,
		AudioFormat:         audioFormatDefault,
		ConcurrentFragments: concurrentFragments,
	}
	downloadDir := dlDir
	fileName := ""
	includeTitleInFilename := false

	// If no options provided, return defaults
	if options == nil {
		return dlOptions, downloadDir, fileName, includeTitleInFilename
	}

	if dlOptions.ConcurrentFragments == 0 || dlOptions.ConcurrentFragments > 20 {
		dlOptions.ConcurrentFragments = consts.ConcurrentFragmentsDefault
	}

	// Override format type if provided
	if options.FormatType != dtypes.FormatTypeNone {
		dlOptions.FormatType = options.FormatType
	}

	// Override video format if provided
	if options.VideoFormat != nil && *options.VideoFormat != dtypes.VideoFormatNone {
		dlOptions.VideoFormat = *options.VideoFormat
	}

	// Override video codec if provided
	if options.VideoCodec != nil && *options.VideoCodec != dtypes.VideoCodecNone {
		dlOptions.VideoCodec = *options.VideoCodec
	}

	// Override video codec if provided
	if options.VideoResolution != nil && *options.VideoResolution != dtypes.VideoResolutionNone {
		dlOptions.VideoResolution = *options.VideoResolution
	}

	// Override audio format if provided
	if options.AudioFormat != nil && *options.AudioFormat != dtypes.AudioFormatNone {
		dlOptions.AudioFormat = *options.AudioFormat
	}

	// Override file name if provided
	if options.Filename != nil {
		fileName = nfile.FileNameWithoutExt(*options.Filename)
	}

	// Include title in filename
	includeTitleInFilename = options.IncludeTitleInFilename

	// Override download directory if provided
	if options.DownloadsDir != nil {
		downloadDir = *options.DownloadsDir
	}

	return dlOptions, downloadDir, fileName, includeTitleInFilename
}
