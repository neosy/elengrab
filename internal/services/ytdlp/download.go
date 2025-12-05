package ytdlpsrv

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"

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
	formatTypeDefault      = dtypes.FormatTypeVideoAudio
	videoFormatDefault     = dtypes.VideoFormatMP4Orig
	audioFormatDefault     = dtypes.AudioFormatMP3
	audioQualityMP3Default = "0"
	audioQualityM4ADefault = "0"
)

func (srv *YtDlpService) Download(url string, options *dservices.DownloadOptions) (<-chan *ddownload.DownloadResult, error) {
	resultCh := make(chan *ddownload.DownloadResult)

	sendResultError := func(err error) {
		resultCh <- &ddownload.DownloadResult{Error: err}
	}

	go func() {
		defer close(resultCh)

		// Prepare download options with defaults and user overrides
		formatType, videoFormat, audioFormat, downloadDir, fileName, includeTitleInFilename := srv.prepareDownloadOptions(options)

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
		args, fileExt, info, err := srv.buildDownloadArgs(url, formatType, videoFormat, audioFormat, srv.options.ConcurrentFragments)
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
	audioFormat dtypes.AudioFormat,
	downloadDir, fileName string,
	includeTitleInFilename bool,
) {
	// Set default values
	formatType = formatTypeDefault
	videoFormat = videoFormatDefault
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

	return formatType, videoFormat, audioFormat, downloadDir, fileName, includeTitleInFilename
}

// buildDownloadArgs build yt-dlp arguments and get file extension and title
func (srv *YtDlpService) buildDownloadArgs(
	url string,
	formatType dtypes.FormatType,
	videoFormat dtypes.VideoFormat,
	audioFormat dtypes.AudioFormat,
	concurrentAragments uint8,
) (args []string, fileExt string, info *dyoutubeinfo.YouTubeInfo, err error) {

	// Default audio quality (used when extracting audio)
	var (
		audioQualityM4A = audioQualityM4ADefault
		audioQualityMP3 = audioQualityMP3Default
	)

	// prevent downloading the entire playlist, only fetch single video
	args = append(args, "--no-playlist")
	args = append(args, "--concurrent-fragments", strconv.Itoa(int(concurrentAragments)))

	switch formatType {
	// Video + Audio or Video only
	case dtypes.FormatTypeVideoAudio, dtypes.FormatTypeVideoOnly:
		var format string

		// Choose video format string based on requested video format
		switch videoFormat {
		case dtypes.VideoFormatOrig:
			format = "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"
		default:
			format = "bestvideo[ext=mp4]+bestaudio[ext=m4a]/bestvideo+bestaudio/best"
		}

		// Get information about the best format from yt-dlp
		var err error
		info, err = srv.getBestFormat(url, format)
		if err != nil {
			return nil, "", nil, err
		}

		// Add format option to yt-dlp arguments
		args = append(args, "-f", format)

		// args = append(args, "--merge-output-format", info.Formats[0].FileExt)

		// Determine output file extension
		switch videoFormat {
		case dtypes.VideoFormatOrig:
			fileExt = info.Formats[0].FileExt
		case dtypes.VideoFormatMP4H264:
			// Codec options (lower - worse)
			// ffmpeg:-c:v libx264 -crf 18 -preset slow -c:a aac -b:a 192k
			// ffmpeg:-c:v libx264 -crf 22 -preset slow -c:a aac -b:a 160k
			// ffmpeg:-c:v libx264 -crf 24 -preset slow -c:a aac -b:a 128k
			args = append(args, "--recode-video", "mp4", "--postprocessor-args", "ffmpeg:-c:v libx264 -crf 22 -preset slow -c:a aac -b:a 160k")
			fileExt = "mp4"
		case dtypes.VideoFormatMP4H265:
			args = append(args, "--recode-video", "mp4", "--postprocessor-args", "ffmpeg:-c:v libx265 -crf 22 -preset slow -c:a aac -b:a 160k")
			fileExt = "mp4"
		default:
			args = append(args, "--recode-video", "mp4")
			fileExt = "mp4"
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
			return nil, "", nil, err
		}

		args = append(args, "-f", format)

		// Choose audio processing based on requested audio format
		switch audioFormat {
		case dtypes.AudioFormatOrig:
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
		default:
			args = append(args, "--extract-audio", "--audio-format", "mp3", "--audio-quality", audioQualityMP3)
			fileExt = "mp3"
		}
	}

	return args, fileExt, info, nil
}
