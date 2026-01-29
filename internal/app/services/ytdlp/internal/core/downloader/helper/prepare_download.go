package helper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/pkg/fnx"
)

// PrepareDownload builds yt-dlp download arguments based on the provided options
func PrepareDownload(
	ctx context.Context,
	url string,
	dlOptions dto.DLOptions,
	bestFormat func(ctx context.Context, url string, format string) (*idto.MediaInfo, error),
) (args []string, fileExt string, dtoMediaInfo *idto.MediaInfo, mediaInfo *ddownload.MediaInfo, err error) {
	var (
		infoVideoCodec = dtypes.VideoCodecNone
		infoAudioCodec = dtypes.AudioCodecNone
	)

	// Default audio quality (used when extracting audio)
	var (
		audioQualityM4A  = audioQualityM4ADefault
		audioQualityMP3  = audioQualityMP3Default
		audioQualityOPUS = audioQualityOPUSDefault
	)

	var mediaFormat idto.MediaFormat

	// prevent downloading the entire playlist, only fetch single video
	args = append(args, "--no-playlist")
	args = append(args, "--no-warnings")

	switch dlOptions.FormatType {
	// Video + Audio or Video only
	case dtypes.FormatTypeVideoAudio, dtypes.FormatTypeVideoOnly:
		var formatQuery, bestFormatQuery string

		var resolution string
		if dlOptions.VideoResolution.Height() > 0 {
			resolution = fmt.Sprintf("[height<=%d][width<=%d]", dlOptions.VideoResolution.Width(), dlOptions.VideoResolution.Width())
		}

		var (
			formatAVC1Query = fmt.Sprintf("bestvideo[ext=mp4][vcodec^=avc1]%s+bestaudio[ext=m4a]", resolution)
			formatAV01Query = fmt.Sprintf("bestvideo[vcodec^=av01]%s+bestaudio[ext=webm]", resolution)
		)

		// Choose video format string based on requested video format
		switch dlOptions.VideoFormat {
		case dtypes.VideoFormatAuto:
			formatQuery = fmt.Sprintf(
				"bestvideo[ext=mp4]%s+bestaudio[ext=m4a]/best[ext=mp4]%s",
				resolution, resolution,
			) + "/best" + fnx.Ternary(resolution != "", resolution, "")
			if dlOptions.VideoCodec == dtypes.VideoCodecH264 {
				formatQuery = fmt.Sprintf("%s/%s", formatAVC1Query, formatQuery)
			}
			if dlOptions.VideoCodec == dtypes.VideoCodecAV1 {
				formatQuery = fmt.Sprintf("%s/%s", formatAV01Query, formatQuery)
			}
			bestFormatQuery = formatQuery
		case dtypes.VideoFormatWebM:
			formatQuery = fmt.Sprintf("bestvideo[ext=webm]%s+bestaudio[ext=webm]", resolution)
			bestFormatQuery = fmt.Sprintf(
				"bestvideo[ext=webm]%s+bestaudio[ext=webm]/bestvideo[ext=webm]%s/best[ext=mp4]%s",
				resolution, resolution, resolution,
			) + "/best" + fnx.Ternary(resolution != "", resolution, "")
			if dlOptions.VideoCodec == dtypes.VideoCodecAV1 {
				formatQuery = fmt.Sprintf("%s/%s", formatAV01Query, formatQuery)
				bestFormatQuery = fmt.Sprintf("%s/%s", formatAV01Query, bestFormatQuery)
			}
		default:
			formatQuery = fmt.Sprintf(
				"bestvideo[ext=mp4]%s+bestaudio[ext=m4a]/bestvideo%s+bestaudio",
				resolution, resolution,
			) + "/best" + fnx.Ternary(resolution != "", resolution, "")
			if dlOptions.VideoCodec == dtypes.VideoCodecH264 {
				formatQuery = fmt.Sprintf("%s/%s", formatAVC1Query, formatQuery)
			}
			if dlOptions.VideoCodec == dtypes.VideoCodecAV1 {
				formatQuery = fmt.Sprintf("%s/%s", formatAV01Query, formatQuery)
			}
			bestFormatQuery = formatQuery
		}

		// Get information about the best format from yt-dlp
		var err error
		dtoMediaInfo, err = bestFormat(ctx, url, bestFormatQuery)
		if err != nil {
			return nil, "", nil, nil, err
		}

		if len(dtoMediaInfo.Formats) == 0 {
			return nil, "", nil, nil, errors.New("not found best format")
		}

		mediaFormat = dtoMediaInfo.Formats[0]
		if len(dtoMediaInfo.Formats) == 2 {
			mediaFormat.ACodec = dtoMediaInfo.Formats[1].ACodec
			mediaFormat.Abr = dtoMediaInfo.Formats[1].Abr
			mediaFormat.Asr = dtoMediaInfo.Formats[1].Asr
		}

		// Add format option to yt-dlp arguments
		args = append(args, "-f", formatQuery)

		// Determine output file extension
		switch dlOptions.VideoFormat {
		case dtypes.VideoFormatAuto:
			fileExt = mediaFormat.FileExt
		case dtypes.VideoFormatWebM:
			fileExt = "webm"
		case dtypes.VideoFormatMP4:
			fileExt = "mp4"
		default:
			fileExt = "mp4"
		}

		// Video resolution
		var (
			isVideoScale   bool
			videoScaleArgs string
		)
		if dlOptions.VideoFormat != dtypes.VideoFormatWebM {
			scaleValue := ScaleValue(
				uint16(mediaFormat.Width),
				uint16(mediaFormat.Height),
				dlOptions.VideoResolution,
			)
			if scaleValue != "" {
				videoScaleArgs = fmt.Sprintf("-vf scale=%s", scaleValue)
			}
			isVideoScale = videoScaleArgs != ""
		}

		infoVideoCodec = mediaFormat.VideoCodec()
		infoAudioCodec = mediaFormat.AudioCodec()

		isInfoFormatWebMCodec := infoVideoCodec == dtypes.VideoCodecAV1 || infoVideoCodec == dtypes.VideoCodecVP9
		isInfoFormatMP4Codec := infoVideoCodec == dtypes.VideoCodecAV1 || infoVideoCodec == dtypes.VideoCodecH264

		var ffmpegArgs string

		// Determine the output video codec
		switch dlOptions.VideoCodec {
		case dtypes.VideoCodecBest:
			if isVideoScale {
				audioCodecArgs := "-c:a copy"
				ffmpegArgs = fmt.Sprintf("%s %s", videoScaleArgs, audioCodecArgs)
			}
		case dtypes.VideoCodecAV1:
			infoVideoCodec = dtypes.VideoCodecAV1
			var (
				isVideoArgs    bool
				isAudioArgs    bool
				videoCodecArgs string
				audioCodecArgs string
			)
			if infoVideoCodec != dtypes.VideoCodecAV1 {
				videoCodecArgs = "-c:v libaom-av1 -crf 0 -b:v 0"
				isVideoArgs = true
			}
			if isVideoArgs || isVideoScale {
				audioCodecArgs = "-c:a copy"
				isAudioArgs = true
			}
			if isVideoArgs || isAudioArgs {
				ffmpegArgs = fmt.Sprintf("%s %s %s", videoScaleArgs, videoCodecArgs, audioCodecArgs)
			}
		case dtypes.VideoCodecH264:
			infoVideoCodec = dtypes.VideoCodecH264
			var (
				isVideoArgs    bool
				isAudioArgs    bool
				videoCodecArgs string
				audioCodecArgs string
			)
			// Codec options (lower - better)
			// ffmpeg:-c:v libx264 -crf 18 -preset slow -c:a aac -b:a 192k
			// ffmpeg:-c:v libx264 -crf 22 -preset slow -c:a aac -b:a 160k
			// ffmpeg:-c:v libx264 -crf 24 -preset slow -c:a aac -b:a 128k
			if infoVideoCodec != dtypes.VideoCodecH264 {
				videoCodecArgs = "-c:v libx264 -crf 22 -preset slow"
				isVideoArgs = true
			}
			if isVideoArgs || isVideoScale {
				audioCodecArgs = "-c:a aac -b:a 160k"
			}
			if isVideoArgs || isAudioArgs {
				ffmpegArgs = fmt.Sprintf("%s %s %s", videoScaleArgs, videoCodecArgs, audioCodecArgs)
				isAudioArgs = true
			}
		case dtypes.VideoCodecH265:
			infoVideoCodec = dtypes.VideoCodecH265
			videoCodecArgs := "-c:v libx265 -crf 22 -preset slow"
			audioCodecArgs := "-c:a aac -b:a 160k"
			ffmpegArgs = fmt.Sprintf("%s %s %s", ffmpegArgs, videoCodecArgs, audioCodecArgs)
		default:
			ffmpegArgs = ""
		}

		ffmpegArgs = strings.TrimSpace(ffmpegArgs)

		switch fileExt {
		case "mp4":
			if ffmpegArgs != "" || !isInfoFormatMP4Codec {
				args = append(args, "--recode-video", "mp4")
			} else {
				args = append(args, "--remux-video", "mp4")
			}
		case "webm":
			if ffmpegArgs != "" || !isInfoFormatWebMCodec {
				args = append(args, "--recode-video", "webm")
			} else {
				args = append(args, "--merge-output-format", "webm")
			}
		}

		if ffmpegArgs != "" {
			args = append(args, "--postprocessor-args", fmt.Sprintf("ffmpeg:%s", ffmpegArgs))
		}

	// Audio only
	case dtypes.FormatTypeAudioOnly:
		var format string

		// Choose audio format string based on requested audio format
		switch dlOptions.AudioFormat {
		case dtypes.AudioFormatM4A:
			format = "bestaudio[ext=m4a]/bestaudio/best"
		default:
			format = "bestaudio/best"
		}

		// Get information about the best audio format
		var err error
		dtoMediaInfo, err = bestFormat(ctx, url, format)
		if err != nil {
			return nil, "", nil, nil, err
		}

		if len(dtoMediaInfo.Formats) == 0 {
			return nil, "", nil, nil, errors.New("not found best format")
		}

		mediaFormat = dtoMediaInfo.Formats[0]
		infoAudioCodec = mediaFormat.AudioCodec()

		args = append(args, "-f", format)

		// Choose audio processing based on requested audio format
		switch dlOptions.AudioFormat {
		case dtypes.AudioFormatAuto:
			if mediaFormat.ACodec != "" && mediaFormat.ACodec == "opus" {
				args = append(args, "--extract-audio", "--audio-format", "opus")
				fileExt = "opus"
			} else {
				fileExt = mediaFormat.FileExt
			}
		case dtypes.AudioFormatM4A:
			infoAudioCodec = dtypes.AudioCodecAAC
			args = append(args, "--extract-audio", "--audio-format", "m4a", "--audio-quality", audioQualityM4A)
			fileExt = "m4a"
		case dtypes.AudioFormatFLAC:
			infoAudioCodec = dtypes.AudioCodecFLAC
			mediaFormat.Abr = audioQualityFLACBitrateDefault
			args = append(args, "--extract-audio", "--audio-format", "flac", "--postprocessor-args", "-compression_level 8")
			fileExt = "flac"
		case dtypes.AudioFormatOPUS:
			infoAudioCodec = dtypes.AudioCodecOPUS
			if mediaFormat.ACodec != "" && mediaFormat.ACodec == "opus" {
				args = append(args, "--extract-audio", "--audio-format", "opus")
			} else {
				args = append(args, "--extract-audio", "--audio-format", "opus", "--audio-quality", audioQualityOPUS)
			}
			fileExt = "opus"
		default:
			infoAudioCodec = dtypes.AudioCodecMP3
			mediaFormat.Abr = audioQualityMP3BitrateDefault
			args = append(args, "--extract-audio", "--audio-format", "mp3", "--audio-quality", audioQualityMP3)
			fileExt = "mp3"
		}
	}

	mediaInfo = &ddownload.MediaInfo{
		FormatType: dlOptions.FormatType,
		Format:     dtypes.MapFileExtToFileFormat(fileExt),
		VideoInfo:  nil,
		AudioInfo:  nil,
	}

	if dtoMediaInfo != nil && len(dtoMediaInfo.Formats) > 0 {
		if infoVideoCodec != dtypes.VideoCodecNone {
			mediaInfo.VideoInfo = &ddownload.VideoInfo{
				Codec: infoVideoCodec,
				Resolution: dtypes.ParseVideoResolutionWH(
					uint16(mediaFormat.Width),
					uint16(mediaFormat.Height),
				),
				Bitrate: int(mediaFormat.Vbr),
				Width:   mediaFormat.Width,
				Height:  mediaFormat.Height,
			}
		}
		if infoAudioCodec != dtypes.AudioCodecNone {
			mediaInfo.AudioInfo = &ddownload.AudioInfo{
				Codec:      infoAudioCodec,
				Bitrate:    int(mediaFormat.Abr),
				SampleRate: mediaFormat.Asr,
			}
		}
	}

	return args, fileExt, dtoMediaInfo, mediaInfo, nil
}
