package ytdlpsrv

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/google/uuid"
	iconstants "github.com/neosy/elengrab/internal/constants"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
	"github.com/neosy/elengrab/pkg/nfile"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

const (
	formatTypeDefault       = dtypes.FormatTypeVideoAudio
	videoFormatDefault      = dtypes.VideoFormatMP4
	videoCodecDefault       = dtypes.VideoCodecBest
	videoResolutionDefault  = dtypes.VideoResolutionBest
	audioFormatDefault      = dtypes.AudioFormatMP3
	audioQualityMP3Default  = "2"
	audioQualityM4ADefault  = "2"
	audioQualityOPUSDefault = "160K"
)

func (srv *YtDlpService) Download(url string, options *dservices.DownloadOptions) (<-chan *ddownload.DownloadResult, error) {
	resultCh := make(chan *ddownload.DownloadResult)

	sendResultError := func(err error) {
		resultCh <- &ddownload.DownloadResult{Error: err}
	}

	go func() {
		defer close(resultCh)

		// Prepare download options with defaults and user overrides
		formatType,
			videoFormat,
			videoCodec,
			videoResolution,
			audioFormat,
			downloadDir,
			fileName,
			includeTitleInFilename := srv.prepareDownloadOptions(options)

		var (
			cmd *exec.Cmd
		)

		// Ensure download directory exists
		if err := checkDir(downloadDir); err != nil {
			sendResultError(fmt.Errorf("failed to check directory: %w", err))
			return
		}

		title, err := srv.getTitleFast(url)
		if err == nil && title != "" {
			resultCh <- &ddownload.DownloadResult{
				YoutubeTitle: title,
			}
		}

		// Build yt-dlp arguments and get file extension and title
		args, fileExt, info, mediaInfo, err := srv.buildDownloadArgs(
			url,
			formatType,
			videoFormat,
			videoCodec,
			videoResolution,
			audioFormat,
			srv.options.ConcurrentFragments,
		)
		if err != nil {
			sendResultError(fmt.Errorf("failed to build download arguments: %w", err))
			return
		}

		// If title is empty, fetch it manually
		title = info.Title
		if title == "" {
			var err error
			title, err = srv.GetTitle(url)
			if err != nil {
				sendResultError(fmt.Errorf("failed to get title: %w", err))
				return
			}
		}

		// Generate a unique filename if none provided
		if fileName == "" {
			fileName = uuid.New().String()
		}

		if includeTitleInFilename {
			fileName = fmt.Sprintf("%s_%s", nfile.SanitizeFileName(title), fileName)
		}

		FileFullName := fmt.Sprintf("%s.%s", fileName, fileExt)
		filePath := path.Join(downloadDir, FileFullName)

		var fileSize *int
		if len(info.Formats) == 1 {
			fileSize = info.Formats[0].Filesize
		}

		resultCh <- &ddownload.DownloadResult{
			YoutubeTitle: title,
			FilePath:     filePath,
			Filename:     fileName,
			FileExt:      fileExt,
			FileFullName: FileFullName,
			Filesize:     fileSize,
		}

		// Add output file path to yt-dlp arguments
		args = append(args, "-o", filePath)

		// Add the video URL to arguments
		args = append(args, url)

		// Execute yt-dlp command
		cmd = exec.Command(srv.cmdPath, args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			// Log full stdout + stderr if yt-dlp fails
			sendResultError(fmt.Errorf("%s failed: %w, output: %s", ytDlpName, err, string(out)))
			return
		}

		// Debug log command output
		srv.logger.Debug(
			"Download",
			"url", url,
			"out", string(out),
		)

		// Get the actual file size
		fileInfo, err := os.Stat(filePath)
		if err == nil {
			fileSize = uptr.Int(int(fileInfo.Size()))
		}

		var partialHash *string
		{
			h, err := nfile.HashPartialMedia(
				filePath,
				iconstants.HashPartialBlocks,
				iconstants.HashPartialBlockSize,
			)
			if err == nil && h != "" {
				partialHash = &h
			}
		}

		// Build response struct
		result := &ddownload.DownloadResult{
			YoutubeTitle: title,
			FilePath:     filePath,
			Filename:     fileName,
			FileExt:      fileExt,
			FileFullName: FileFullName,
			Filesize:     fileSize,
			PartialHash:  partialHash,
			MediaInfo:    mediaInfo,
		}

		resultCh <- result

		// Log successful download
		srv.logger.Info("Download successful", "info", result)
	}()

	return resultCh, nil
}

// prepareDownloadOptions prepare download options with defaults and user overrides
func (srv *YtDlpService) prepareDownloadOptions(options *dservices.DownloadOptions) (
	formatType dtypes.FormatType,
	videoFormat dtypes.VideoFormat,
	videoCodec dtypes.VideoCodec,
	videoResolution dtypes.VideoResolution,
	audioFormat dtypes.AudioFormat,
	downloadDir, fileName string,
	includeTitleInFilename bool,
) {
	// Set default values
	formatType = formatTypeDefault
	videoFormat = videoFormatDefault
	videoCodec = videoCodecDefault
	videoResolution = videoResolutionDefault
	audioFormat = audioFormatDefault
	downloadDir = srv.downloadsDir
	fileName = ""
	includeTitleInFilename = false

	// If no options provided, return defaults
	if options == nil {
		return
	}

	// Override format type if provided
	if options.FormatType != dtypes.FormatTypeNone {
		formatType = options.FormatType
	}

	// Override video format if provided
	if options.VideoFormat != nil && *options.VideoFormat != dtypes.VideoFormatNone {
		videoFormat = *options.VideoFormat
	}

	// Override video codec if provided
	if options.VideoCodec != nil && *options.VideoCodec != dtypes.VideoCodecNone {
		videoCodec = *options.VideoCodec
	}

	// Override video codec if provided
	if options.VideoResolution != nil && *options.VideoResolution != dtypes.VideoResolutionNone {
		videoResolution = *options.VideoResolution
	}

	// Override audio format if provided
	if options.AudioFormat != nil && *options.AudioFormat != dtypes.AudioFormatNone {
		audioFormat = *options.AudioFormat
	}

	// Override file name if provided
	if options.Filename != nil {
		fileName = fileNameWithoutExt(*options.Filename)
	}

	// Include title in filename
	includeTitleInFilename = options.IncludeTitleInFilename

	// Override download directory if provided
	if options.DownloadsDir != nil {
		downloadDir = *options.DownloadsDir
	}

	return formatType, videoFormat, videoCodec, videoResolution, audioFormat, downloadDir, fileName, includeTitleInFilename
}

// buildDownloadArgs build yt-dlp arguments and get file extension and title
func (srv *YtDlpService) buildDownloadArgs(
	url string,
	formatType dtypes.FormatType,
	videoFormat dtypes.VideoFormat,
	videoCodec dtypes.VideoCodec,
	videoResolution dtypes.VideoResolution,
	audioFormat dtypes.AudioFormat,
	concurrentAragments uint8,
) (args []string, fileExt string, info *dyoutubeinfo.YouTubeInfo, mediaInfo *ddownload.MediaInfo, err error) {
	var (
		infoVideoCodec = dtypes.VideoCodecNone
	)

	// Default audio quality (used when extracting audio)
	var (
		audioQualityM4A  = audioQualityM4ADefault
		audioQualityMP3  = audioQualityMP3Default
		audioQualityOPUS = audioQualityOPUSDefault
	)

	// prevent downloading the entire playlist, only fetch single video
	args = append(args, "--no-playlist")
	args = append(args, "--concurrent-fragments", strconv.Itoa(int(concurrentAragments)))

	switch formatType {
	// Video + Audio or Video only
	case dtypes.FormatTypeVideoAudio, dtypes.FormatTypeVideoOnly:
		var format, infoFormat string

		var resolution string
		if videoResolution.Height() > 0 {
			resolution = fmt.Sprintf("[height<=%d][width<=%d]", videoResolution.Width(), videoResolution.Width())
		}

		var (
			formatAVC1Query = fmt.Sprintf("bestvideo[ext=mp4][vcodec^=avc1]%s+bestaudio[ext=m4a]", resolution)
			formatAV01Query = fmt.Sprintf("bestvideo[vcodec^=av01]%s+bestaudio[ext=webm]", resolution)
		)

		// Choose video format string based on requested video format
		switch videoFormat {
		case dtypes.VideoFormatBest:
			format = fmt.Sprintf(
				"bestvideo[ext=mp4]%s+bestaudio[ext=m4a]/best[ext=mp4]%s/best%s/best",
				resolution, resolution, resolution,
			)
			if videoCodec == dtypes.VideoCodecH264 {
				format = fmt.Sprintf("%s/%s", formatAVC1Query, format)
			}
			if videoCodec == dtypes.VideoCodecAV1 {
				format = fmt.Sprintf("%s/%s", formatAV01Query, format)
			}
			infoFormat = format
		case dtypes.VideoFormatWebM:
			format = fmt.Sprintf("bestvideo[ext=webm]%s+bestaudio[ext=webm]", resolution)
			infoFormat = fmt.Sprintf(
				"bestvideo[ext=webm]%s+bestaudio[ext=webm]/bestvideo[ext=webm]%s/best[ext=mp4]%s/best%s/best",
				resolution, resolution, resolution, resolution,
			)
			if videoCodec == dtypes.VideoCodecAV1 {
				format = fmt.Sprintf("%s/%s", formatAV01Query, format)
				infoFormat = fmt.Sprintf("%s/%s", formatAV01Query, infoFormat)
			}
		default:
			format = fmt.Sprintf(
				"bestvideo[ext=mp4]%s+bestaudio[ext=m4a]/bestvideo%s+bestaudio/best%s/best",
				resolution, resolution, resolution,
			)
			if videoCodec == dtypes.VideoCodecH264 {
				format = fmt.Sprintf("%s/%s", formatAVC1Query, format)
			}
			if videoCodec == dtypes.VideoCodecAV1 {
				format = fmt.Sprintf("%s/%s", formatAV01Query, format)
			}
			infoFormat = format
		}

		// Get information about the best format from yt-dlp
		var err error
		info, err = srv.getBestFormat(url, infoFormat)
		if err != nil {
			return nil, "", nil, nil, err
		}

		if len(info.Formats) == 0 {
			return nil, "", nil, nil, errors.New("not found best format")
		}

		// Add format option to yt-dlp arguments
		args = append(args, "-f", format)

		// Determine output file extension
		switch videoFormat {
		case dtypes.VideoFormatBest:
			fileExt = info.Formats[0].FileExt
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
		if videoFormat != dtypes.VideoFormatWebM {
			scaleValue := srv.scaleValue(uint16(info.Formats[0].Width), uint16(info.Formats[0].Height), videoResolution)
			if scaleValue != "" {
				videoScaleArgs = fmt.Sprintf("-vf scale=%s", scaleValue)
			}
			isVideoScale = videoScaleArgs != ""
		}

		var (
			isInfoFormatWebMCodec = false
			isInfoFormatMP4Codec  = false
		)
		if info.Formats[0].VCodec != nil {
			isInfoFormatWebMCodec = strings.HasPrefix(*info.Formats[0].VCodec, "av01") ||
				strings.HasPrefix(*info.Formats[0].VCodec, "vp9")
			isInfoFormatMP4Codec = strings.HasPrefix(*info.Formats[0].VCodec, "av01") ||
				strings.HasPrefix(*info.Formats[0].VCodec, "avc1")
			if strings.HasPrefix(*info.Formats[0].VCodec, "av01") {
				infoVideoCodec = dtypes.VideoCodecAV1
			}
			if strings.HasPrefix(*info.Formats[0].VCodec, "vp9") {
				infoVideoCodec = dtypes.VideoCodecVP9
			}
			if strings.HasPrefix(*info.Formats[0].VCodec, "avc1") {
				infoVideoCodec = dtypes.VideoCodecH264
			}
		}

		var ffmpegArgs string

		// Determine the output video codec
		switch videoCodec {
		case dtypes.VideoCodecBest:
			if isVideoScale {
				audioCodecArgs := "-c:a copy"
				ffmpegArgs = fmt.Sprintf("%s %s", videoScaleArgs, audioCodecArgs)
			}
		case dtypes.VideoCodecAV1:
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
		switch audioFormat {
		case dtypes.AudioFormatM4A:
			format = "bestaudio[ext=m4a]/bestaudio"
		default:
			format = "bestaudio"
		}

		// Get information about the best audio format
		var err error
		info, err = srv.getBestFormat(url, format)
		if err != nil {
			return nil, "", nil, nil, err
		}

		args = append(args, "-f", format)

		// Choose audio processing based on requested audio format
		switch audioFormat {
		case dtypes.AudioFormatBest:
			format := info.Formats[0]
			if format.ACodec != nil && *format.ACodec == "opus" {
				args = append(args, "--extract-audio", "--audio-format", "opus")
				fileExt = "opus"
			} else {
				fileExt = format.FileExt
			}
		case dtypes.AudioFormatM4A:
			args = append(args, "--extract-audio", "--audio-format", "m4a", "--audio-quality", audioQualityM4A)
			fileExt = "m4a"
		case dtypes.AudioFormatFLAC:
			args = append(args, "--extract-audio", "--audio-format", "flac", "--postprocessor-args", "-compression_level 8")
			fileExt = "flac"
		case dtypes.AudioFormatOPUS:
			format := info.Formats[0]
			if format.ACodec != nil && *format.ACodec == "opus" {
				args = append(args, "--extract-audio", "--audio-format", "opus")
			} else {
				args = append(args, "--extract-audio", "--audio-format", "opus", "--audio-quality", audioQualityOPUS)
			}
			fileExt = "opus"
		default:
			args = append(args, "--extract-audio", "--audio-format", "mp3", "--audio-quality", audioQualityMP3)
			fileExt = "mp3"
		}
	}

	mediaInfo = &ddownload.MediaInfo{
		FormatType: formatType,
		Format:     dtypes.MapFileExtToFileFormat(fileExt),
		VideoCodec: videoCodec,
		Resolution: videoResolution,
		Width:      int(videoResolution.Width()),
		Height:     int(videoResolution.Height()),
	}

	if info != nil && len(info.Formats) > 0 {
		infoFormat := info.Formats[0]
		if mediaInfo.VideoCodec == dtypes.VideoCodecBest {
			mediaInfo.VideoCodec = infoVideoCodec
		}
		if mediaInfo.Resolution == dtypes.VideoResolutionBest {
			mediaInfo.Width = infoFormat.Width
			mediaInfo.Height = infoFormat.Height
			mediaInfo.Resolution = dtypes.ParseVideoResolutionWH(uint16(mediaInfo.Width), uint16(mediaInfo.Height))
		}
	}

	return args, fileExt, info, mediaInfo, nil
}

// scaleValue builds an ffmpeg scale expression based on source dimensions
// and target resolution, preserving aspect ratio.
//
// width, height — source dimensions (may be rotated)
// toResolution  — target resolution (Width/Height)
func (srv *YtDlpService) scaleValue(width, height uint16, toResolution dtypes.VideoResolution) string {
	// Target dimensions (e.g. 1280x720)
	targetW := toResolution.Width()
	targetH := toResolution.Height()

	// Determine source orientation:
	// true  -> landscape (width >= height)
	// false -> portrait (height > width)
	sourceIsLandscape := width >= height

	// Normalize source dimensions so srcW >= srcH for comparison
	srcW, srcH := width, height
	if !sourceIsLandscape {
		// swap for portrait to compare logically as W/H
		srcW, srcH = height, width
	}

	// If any dimension is zero, return empty (no scale)
	if targetW == 0 || targetH == 0 || srcW == 0 || srcH == 0 {
		return ""
	}

	// If source is already within target bounds, no scaling needed
	if srcW <= targetW || srcH <= targetH {
		return ""
	}

	// Build scale value for ffmpeg scale filter:
	// - landscape: limit height (targetH) -> "-1:targetH"
	// - portrait:  limit width  (targetH per your convention) -> "targetH:-1"
	// Note: we use -1 so ffmpeg recalculates the other dimension preserving aspect ratio.
	var scaleValue string
	if sourceIsLandscape {
		// landscape -> fix height
		scaleValue = fmt.Sprintf("-1:%d", targetH)
	} else {
		// portrait  -> fix width (using targetH per your convention)
		scaleValue = fmt.Sprintf("%d:-1", targetH)
	}

	return scaleValue
}
