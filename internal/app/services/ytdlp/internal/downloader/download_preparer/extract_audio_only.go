package downloadpreparer

import (
	"context"
	"errors"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (p *DownloadPreparer) extractAudioOnly(
	ctx context.Context,
	url string,
	dlOptions idto.DLOptions,
) (extractResult, error) {
	var formatQuery string

	// Choose audio format string based on requested audio format
	switch dlOptions.AudioFormat {
	case dtypes.AudioFormatM4A:
		formatQuery = "bestaudio[ext=m4a]/bestaudio/best"
	default:
		formatQuery = "bestaudio/best"
	}

	// Get information about the best audio format
	var err error
	extractInfo, err := p.getExtractInfo(
		ctx,
		url,
		formatQuery,
		idto.WithUseCookies(dlOptions.CookieFilePathIfNeeded()),
		idto.WithEnsureCache(false),
	)
	if err != nil {
		return extractResult{}, err
	}

	if len(extractInfo.Formats) == 0 {
		return extractResult{}, errors.New("not found best format")
	}

	mediaFormat := extractInfo.Formats[0]

	selectedFormatIDs := append([]string{}, mediaFormat.FormatID)

	srcAudioFormat := dtypes.MapFileExtToFileFormat(mediaFormat.FileExt).AudioFormat()
	srcAudioCodec := mediaFormat.AudioCodec()

	outAudioFormat := dlOptions.AudioFormat
	outAudioCodec := srcAudioCodec

	args := make([]string, 0)

	if outAudioFormat == dtypes.AudioFormatAuto {
		if srcAudioFormat == dtypes.AudioFormatNone {
			outAudioFormat = outAudioCodec.AudioFormat()
		}
	}
	if outAudioFormat == dtypes.AudioFormatNone {
		outAudioFormat = dtypes.AudioFormatMP3
	}

	// Choose audio processing based on requested audio format
	switch outAudioFormat {
	case dtypes.AudioFormatAuto:
		if srcAudioCodec == dtypes.AudioCodecOPUS {
			args = append(args, "--extract-audio", "--audio-format", "opus")
			outAudioFormat = dtypes.AudioFormatOPUS
		} else {
			outAudioFormat = srcAudioFormat
		}
	case dtypes.AudioFormatM4A:
		outAudioCodec = dtypes.AudioCodecAAC
		args = append(args, "--extract-audio", "--audio-format", "m4a", "--audio-quality", audioQualityM4ADefault)
	case dtypes.AudioFormatFLAC:
		outAudioCodec = dtypes.AudioCodecFLAC
		mediaFormat.Abr = audioQualityFLACBitrateDefault
		args = append(args, "--extract-audio", "--audio-format", "flac", "--postprocessor-args", "-compression_level 8")
	case dtypes.AudioFormatOPUS:
		outAudioCodec = dtypes.AudioCodecOPUS
		if mediaFormat.ACodec != "" && mediaFormat.ACodec == "opus" {
			args = append(args, "--extract-audio", "--audio-format", "opus")
		} else {
			args = append(args, "--extract-audio", "--audio-format", "opus", "--audio-quality", audioQualityOPUSDefault)
		}
	default:
		outAudioFormat = dtypes.AudioFormatMP3
		outAudioCodec = dtypes.AudioCodecMP3
		mediaFormat.Abr = audioQualityMP3BitrateDefault
		args = append(args, "--extract-audio", "--audio-format", "mp3", "--audio-quality", audioQualityMP3Default)
	}

	return extractResult{
		formatQuery: formatQuery,
		args:        args,

		fileExt: outAudioFormat.FileFormat().Ext(),

		info:        extractInfo,
		mediaFormat: &mediaFormat,

		audioCodec: outAudioCodec,

		selectedFormatIDs: selectedFormatIDs,
	}, nil
}
