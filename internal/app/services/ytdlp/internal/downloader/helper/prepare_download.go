package helper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/executor"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/ffmpeg"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/pkg/fnx"
)

// PrepareDownload builds yt-dlp download arguments based on the provided options
func PrepareDownload(
	ctx context.Context,
	url string,
	dlOptions idto.DLOptions,
	bestFormat func(ctx context.Context, url string, format string, opts ...executor.Option) (*idto.MediaInfo, error),
) (args []string, fileExt string, dtoMediaInfo *idto.MediaInfo, mediaInfo *ddownload.MediaInfo, err error) {
	var (
		outVideoCodec = dtypes.VideoCodecNone
		outAudioCodec = dtypes.AudioCodecNone
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
			formatAVC1Query = fmt.Sprintf(
				"bestvideo[ext=mp4][vcodec^=avc1]%s+bestaudio[ext=m4a]/bestvideo[vcodec^=avc1]%s+bestaudio",
				resolution, resolution,
			)
			formatAV01Query = fmt.Sprintf("bestvideo[vcodec^=av01]%s+bestaudio/bestvideo[vcodec^=av01]+bestaudio", resolution)
		)

		// Choose video format string based on requested video format
		switch dlOptions.VideoFormat {
		case dtypes.VideoFormatAuto:
			formatQuery = fmt.Sprintf(
				"bestvideo[ext=webm]%s+bestaudio[ext=webm]/bestvideo[ext=mp4]%s+bestaudio[ext=m4a]/best%s",
				resolution, resolution, resolution,
			) + fnx.Ternary(resolution != "", "bestvideo[ext=webm]+bestaudio[ext=webm]/bestvideo[ext=mp4]+bestaudio[ext=m4a]/best", "")
			if dlOptions.VideoCodec == dtypes.VideoCodecH264 {
				formatQuery = fmt.Sprintf("%s/%s", formatAVC1Query, formatQuery)
			}
			if dlOptions.VideoCodec == dtypes.VideoCodecAV1 {
				formatQuery = fmt.Sprintf("%s/%s", formatAV01Query, formatQuery)
			}
			bestFormatQuery = formatQuery
		case dtypes.VideoFormatWebM:
			formatQuery = fmt.Sprintf(
				"bestvideo[ext=webm]%s+bestaudio/bestvideo[vcodec^=av01]%s+bestaudio/bestvideo[vcodec^=vp09]%s+bestaudio/bestvideo%s+bestaudio",
				resolution, resolution, resolution, resolution,
			) + fnx.Ternary(resolution != "", "/bestvideo[vcodec^=av01]+bestaudio/bestvideo[vcodec^=vp09]+bestaudio/bestvideo+bestaudio", "") +
				fnx.Ternary(resolution != "", fmt.Sprintf("/best%s/best", resolution), "/best")
			if dlOptions.VideoCodec == dtypes.VideoCodecAV1 {
				formatQuery = fmt.Sprintf("%s/%s", formatAV01Query, formatQuery)
			}
			bestFormatQuery = formatQuery
		default:
			formatQuery = fmt.Sprintf(
				"bestvideo[ext=mp4]%s+bestaudio[ext=m4a]/bestvideo%s+bestaudio",
				resolution, resolution,
			) + fnx.Ternary(resolution != "", "/bestvideo+bestaudio/best", "/best")
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
		dtoMediaInfo, err = bestFormat(
			ctx,
			url,
			bestFormatQuery,
			executor.WithUseCookies(dlOptions.RequiresYouTubeCookies),
			executor.WithEnsureCache(false),
		)
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

		srcVideoFormat := dtypes.MapFileExtToFileFormat(mediaFormat.FileExt).VideoFormat()
		srcVideoCodec := mediaFormat.VideoCodec()

		outVideoFormat := dlOptions.VideoFormat
		outVideoCodec = mediaFormat.VideoCodec()

		outAudioCodec = mediaFormat.AudioCodec()

		isSrcFormatWebMCodec := srcVideoCodec == dtypes.VideoCodecAV1 || srcVideoCodec == dtypes.VideoCodecVP9
		isSrcFormatMP4Codec := srcVideoCodec == dtypes.VideoCodecAV1 || srcVideoCodec == dtypes.VideoCodecH264

		// Determine output video format
		switch outVideoFormat {
		case dtypes.VideoFormatAuto:
			outVideoFormat = dlOptions.VideoCodec.Format()
			// fileExt = dlOptions.VideoCodec.Format().FileFormat().Ext()
			if outVideoFormat == dtypes.VideoFormatNone && isSrcFormatWebMCodec {
				outVideoFormat = dtypes.VideoFormatWebM
			}
			if outVideoFormat == dtypes.VideoFormatNone {
				outVideoFormat = srcVideoFormat
			}
		case dtypes.VideoFormatWebM:
		case dtypes.VideoFormatMP4:
		default:
			outVideoFormat = dtypes.VideoFormatMP4
		}

		fileExt = outVideoFormat.FileFormat().Ext()

		// Video resolution
		var (
			isVideoScale   bool
			videoScaleArgs string
		)
		{
			resolution := dtypes.ParseVideoResolutionWH(uint16(mediaFormat.Width), uint16(mediaFormat.Height))
			scaleValue := ScaleValue(
				resolution.Width(),
				resolution.Height(),
				dlOptions.VideoResolution,
			)
			if scaleValue != "" {
				videoScaleArgs = fmt.Sprintf("-vf scale=%s", scaleValue)
				isVideoScale = true
			}
		}

		var ffmpegArgs string

		downloadVideoCodec := dlOptions.VideoCodec
		if outVideoFormat == dtypes.VideoFormatWebM &&
			dlOptions.VideoCodec == dtypes.VideoCodecBest &&
			!isSrcFormatWebMCodec {
			downloadVideoCodec = dtypes.VideoCodecAV1
		}

		// Determine the output video codec
		switch downloadVideoCodec {
		case dtypes.VideoCodecBest:
			needConvert := isVideoScale ||
				(outVideoFormat == dtypes.VideoFormatWebM && outAudioCodec != dtypes.AudioCodecOPUS)
			if needConvert {
				videoCodecArgs := "copy"
				if videoScaleArgs != "" {
					videoCodecArgs = ffmpeg.VideoEncoderArgs(outVideoCodec)
				}
				audioCodecArgs := "copy"
				if outVideoFormat == dtypes.VideoFormatWebM {
					audioCodecArgs = ffmpeg.AudioEncoderArgs(dtypes.AudioCodecOPUS)
				}
				ffmpegArgs = fmt.Sprintf(
					"%s -c:v %s -c:a %s",
					videoScaleArgs,
					videoCodecArgs,
					audioCodecArgs,
				)
			}
		case dtypes.VideoCodecAV1:
			outVideoCodec = dtypes.VideoCodecAV1
			needConvert := isVideoScale ||
				srcVideoCodec != dtypes.VideoCodecAV1 ||
				(outVideoFormat == dtypes.VideoFormatWebM && outAudioCodec != dtypes.AudioCodecOPUS)
			if needConvert {
				videoCodecArgs := "copy"
				if isVideoScale || srcVideoCodec != dtypes.VideoCodecAV1 {
					videoCodecArgs = ffmpeg.VideoEncoderArgs(outVideoCodec)
				}
				audioCodecArgs := "copy"
				if outVideoFormat == dtypes.VideoFormatWebM {
					audioCodecArgs = ffmpeg.AudioEncoderArgs(dtypes.AudioCodecOPUS)
				}
				ffmpegArgs = fmt.Sprintf(
					"%s -c:v %s -c:a %s",
					videoScaleArgs,
					videoCodecArgs,
					audioCodecArgs,
				)
			}
		case dtypes.VideoCodecH264:
			// Codec options (lower - better)
			// ffmpeg:-c:v libx264 -crf 18 -preset slow -c:a aac -b:a 192k
			// ffmpeg:-c:v libx264 -crf 22 -preset slow -c:a aac -b:a 160k
			// ffmpeg:-c:v libx264 -crf 24 -preset slow -c:a aac -b:a 128k
			outVideoCodec = dtypes.VideoCodecH264
			needConvert := isVideoScale || outVideoCodec != srcVideoCodec
			if needConvert {
				videoCodecArgs := ffmpeg.VideoEncoderArgs(outVideoCodec)
				audioCodecArgs := ffmpeg.AudioEncoderArgs(dtypes.AudioCodecAAC)
				ffmpegArgs = fmt.Sprintf(
					"%s -c:v %s -c:a %s",
					videoScaleArgs,
					videoCodecArgs,
					audioCodecArgs,
				)
			}
		case dtypes.VideoCodecH265:
			outVideoCodec = dtypes.VideoCodecH265
			needConvert := isVideoScale || outVideoCodec != srcVideoCodec
			if needConvert {
				videoCodecArgs := ffmpeg.VideoEncoderArgs(outVideoCodec)
				audioCodecArgs := ffmpeg.AudioEncoderArgs(dtypes.AudioCodecAAC)
				ffmpegArgs = fmt.Sprintf(
					"%s -c:v %s -c:a %s",
					videoScaleArgs,
					videoCodecArgs,
					audioCodecArgs,
				)
			}
		default:
			ffmpegArgs = ""
		}

		ffmpegArgs = strings.TrimSpace(ffmpegArgs)

		switch outVideoFormat {
		case dtypes.VideoFormatMP4:
			if ffmpegArgs != "" || !isSrcFormatMP4Codec {
				args = append(args, "--recode-video", "mp4")
			} else {
				args = append(args, "--remux-video", "mp4")
			}
		case dtypes.VideoFormatWebM:
			if ffmpegArgs != "" || !isSrcFormatWebMCodec || outVideoFormat != srcVideoFormat {
				args = append(args, "--recode-video", "webm")
			} else {
				args = append(args, "--merge-output-format", "webm")
			}
		}

		if ffmpegArgs != "" {
			if srcVideoFormat == outVideoFormat {
				args = append(args, "--ppa", fmt.Sprintf("ffmpeg:%s", ffmpegArgs))
			} else {
				args = append(args, "--ppa", fmt.Sprintf("VideoConvertor:%s", ffmpegArgs))
			}
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
		dtoMediaInfo, err = bestFormat(
			ctx,
			url,
			format,
			executor.WithUseCookies(dlOptions.RequiresYouTubeCookies),
			executor.WithEnsureCache(false),
		)
		if err != nil {
			return nil, "", nil, nil, err
		}

		if len(dtoMediaInfo.Formats) == 0 {
			return nil, "", nil, nil, errors.New("not found best format")
		}

		mediaFormat = dtoMediaInfo.Formats[0]

		srcAudioFormat := dtypes.MapFileExtToFileFormat(mediaFormat.FileExt).AudioFormat()
		srcAudioCodec := mediaFormat.AudioCodec()

		outAudioFormat := dlOptions.AudioFormat
		outAudioCodec = srcAudioCodec

		args = append(args, "-f", format)

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
			args = append(args, "--extract-audio", "--audio-format", "m4a", "--audio-quality", audioQualityM4A)
		case dtypes.AudioFormatFLAC:
			outAudioCodec = dtypes.AudioCodecFLAC
			mediaFormat.Abr = audioQualityFLACBitrateDefault
			args = append(args, "--extract-audio", "--audio-format", "flac", "--postprocessor-args", "-compression_level 8")
		case dtypes.AudioFormatOPUS:
			outAudioCodec = dtypes.AudioCodecOPUS
			if mediaFormat.ACodec != "" && mediaFormat.ACodec == "opus" {
				args = append(args, "--extract-audio", "--audio-format", "opus")
			} else {
				args = append(args, "--extract-audio", "--audio-format", "opus", "--audio-quality", audioQualityOPUS)
			}
		default:
			outAudioFormat = dtypes.AudioFormatMP3
			outAudioCodec = dtypes.AudioCodecMP3
			mediaFormat.Abr = audioQualityMP3BitrateDefault
			args = append(args, "--extract-audio", "--audio-format", "mp3", "--audio-quality", audioQualityMP3)
		}
		fileExt = outAudioFormat.FileFormat().Ext()
	}

	mediaInfo = &ddownload.MediaInfo{
		FormatType: dlOptions.FormatType,
		Format:     dtypes.MapFileExtToFileFormat(fileExt),
		VideoInfo:  nil,
		AudioInfo:  nil,
	}

	if dtoMediaInfo != nil && len(dtoMediaInfo.Formats) > 0 {
		if outVideoCodec != dtypes.VideoCodecNone {
			mediaInfo.VideoInfo = &ddownload.VideoInfo{
				Codec: outVideoCodec,
				Resolution: dtypes.ParseVideoResolutionWH(
					uint16(mediaFormat.Width),
					uint16(mediaFormat.Height),
				),
				Bitrate: int(mediaFormat.Vbr),
				Width:   int(mediaFormat.Width),
				Height:  int(mediaFormat.Height),
			}
		}
		if outAudioCodec != dtypes.AudioCodecNone {
			mediaInfo.AudioInfo = &ddownload.AudioInfo{
				Codec:      outAudioCodec,
				Bitrate:    int(mediaFormat.Abr),
				SampleRate: mediaFormat.Asr,
			}
		}
	}

	return args, fileExt, dtoMediaInfo, mediaInfo, nil
}
