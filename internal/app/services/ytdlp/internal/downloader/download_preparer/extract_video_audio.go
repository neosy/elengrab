package downloadpreparer

import (
	"context"
	"fmt"
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/fnx"
)

func (p *DownloadPreparer) extractVideoAudio(
	ctx context.Context,
	url string,
	dlOptions idto.DLOptions,
) (extractResult, error) {
	formatQuery := p.buildVideoFormatQuery(dlOptions)

	fetchedInfo, err := p.fetchMediaInfo(
		ctx,
		url,
		formatQuery,
		dlOptions,
	)
	if err != nil {
		return extractResult{}, err
	}

	downloadParams := p.buildVideoDownloadParams(
		fetchedInfo.mediaFormat,
		dlOptions,
	)

	args := make([]string, 0)
	args = append(args, downloadParams.args...)

	return extractResult{
		formatQuery: formatQuery,
		args:        args,

		fileExt: downloadParams.fileExt,

		info:        fetchedInfo.info,
		mediaFormat: fetchedInfo.mediaFormat,

		videoCodec: downloadParams.videoCodec,
		audioCodec: downloadParams.audioCodec,

		selectedFormatIDs: fetchedInfo.selectedFormatIDs,
	}, nil
}

func (p *DownloadPreparer) buildVideoFormatQuery(dlOptions idto.DLOptions) string {
	var (
		resolution  string
		formatQuery string
	)

	if dlOptions.VideoResolution.Height() > 0 {
		// TODO: improve resolution query to better match requested resolution
		// var (
		// 	prevResolution = dlOptions.VideoResolution.Prev()
		// 	prevWidth      = prevResolution.Width()
		// 	prevHeight     = prevResolution.Height()
		// )

		// resolution = fmt.Sprintf(
		// 	"[height>%d][height<=%d][width>%d][width<=%d]/[height>%d][height<=%d][width>%d][width<=%d]/[height<=%d][width<=%d]",
		// 	prevHeight, dlOptions.VideoResolution.Height(), prevWidth, dlOptions.VideoResolution.Width(),
		// 	prevWidth, dlOptions.VideoResolution.Width(), prevHeight, dlOptions.VideoResolution.Height(),
		// 	dlOptions.VideoResolution.Width(), dlOptions.VideoResolution.Width(),
		// )
		resolution = fmt.Sprintf(
			"[height<=%d][width<=%d]",
			dlOptions.VideoResolution.Width(), dlOptions.VideoResolution.Width(),
		)
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
	case dtypes.VideoFormatWebM:
		formatQuery = fmt.Sprintf(
			"bestvideo[ext=webm]%s+bestaudio/bestvideo[vcodec^=av01]%s+bestaudio/bestvideo[vcodec^=vp09]%s+bestaudio/bestvideo%s+bestaudio",
			resolution, resolution, resolution, resolution,
		) + fnx.Ternary(resolution != "", "/bestvideo[vcodec^=av01]+bestaudio/bestvideo[vcodec^=vp09]+bestaudio/bestvideo+bestaudio", "") +
			fnx.Ternary(resolution != "", fmt.Sprintf("/best%s/best", resolution), "/best")
		if dlOptions.VideoCodec == dtypes.VideoCodecAV1 {
			formatQuery = fmt.Sprintf("%s/%s", formatAV01Query, formatQuery)
		}
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
	}

	return formatQuery
}

type videoDownloadParams struct {
	fileExt string
	args    []string

	videoCodec dtypes.VideoCodec
	audioCodec dtypes.AudioCodec
}

func (p *DownloadPreparer) buildVideoDownloadParams(
	mediaFormat *idto.ExtractMediaFormat,
	dlOptions idto.DLOptions,
) videoDownloadParams {
	srcVideoFormat := dtypes.MapFileExtToFileFormat(mediaFormat.FileExt).VideoFormat()
	srcVideoCodec := mediaFormat.VideoCodec()

	outVideoFormat := dlOptions.VideoFormat
	outVideoCodec := mediaFormat.VideoCodec()

	outAudioCodec := mediaFormat.AudioCodec()

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

	// Video resolution
	var (
		isVideoScale   bool
		videoScaleArgs string
	)
	{
		resolution := dtypes.ParseVideoResolutionWH(uint16(mediaFormat.Width), uint16(mediaFormat.Height))
		scaleValue := helper.ScaleValue(
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
				videoCodecArgs = helper.VideoEncoderArgs(outVideoCodec)
			}
			audioCodecArgs := "copy"
			if outVideoFormat == dtypes.VideoFormatWebM {
				audioCodecArgs = helper.AudioEncoderArgs(dtypes.AudioCodecOPUS)
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
				videoCodecArgs = helper.VideoEncoderArgs(outVideoCodec)
			}
			audioCodecArgs := "copy"
			if outVideoFormat == dtypes.VideoFormatWebM {
				audioCodecArgs = helper.AudioEncoderArgs(dtypes.AudioCodecOPUS)
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
			videoCodecArgs := helper.VideoEncoderArgs(outVideoCodec)
			audioCodecArgs := helper.AudioEncoderArgs(dtypes.AudioCodecAAC)
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
			videoCodecArgs := helper.VideoEncoderArgs(outVideoCodec)
			audioCodecArgs := helper.AudioEncoderArgs(dtypes.AudioCodecAAC)
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

	args := make([]string, 0)

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

	return videoDownloadParams{
		fileExt: outVideoFormat.FileFormat().Ext(),
		args:    args,

		videoCodec: outVideoCodec,
		audioCodec: outAudioCodec,
	}
}
