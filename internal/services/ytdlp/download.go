package ytdlpsrv

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

const (
	formatTypeDefault      = dtypes.FormatTypeVideoAudio
	videoFormatDefault     = dtypes.VideoFormatMP4
	audioFormatDefault     = dtypes.AudioFormatMP3
	audioQualityMP3Default = "0"
	audioQualityM4ADefault = "0"
)

func (srv *YtDlpService) Download(url string, options *ddownload.DownloadOptions) (<-chan *ddownload.DownloadResult, error) {
	resultCh := make(chan *ddownload.DownloadResult)

	sendResultError := func(err error) {
		resultCh <- &ddownload.DownloadResult{Error: err}
	}

	go func() {
		defer close(resultCh)

		// Prepare download options with defaults and user overrides
		formatType, videoFormat, audioFormat, downloadDir, fileName := srv.prepareDownloadOptions(options)

		var (
			cmd *exec.Cmd
		)

		// Ensure download directory exists
		if err := checkDir(downloadDir); err != nil {
			srv.logger.Error(err.Error())
			sendResultError(err)
			return
		}

		// Build yt-dlp arguments and get file extension and title
		args, fileExt, info, err := srv.buildDownloadArgs(url, formatType, videoFormat, audioFormat)
		if err != nil {
			srv.logger.Error(err.Error())
			sendResultError(err)
			return
		}

		// If title is empty, fetch it manually
		title := info.Title
		if title == "" {
			var err error
			title, err = srv.GetTitle(url)
			if err != nil {
				srv.logger.Error(err.Error())
				sendResultError(err)
				return
			}
		}

		// Generate a unique filename if none provided
		if fileName == "" {
			fileName = uuid.New().String()
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
		outByte, err := cmd.CombinedOutput()
		if err != nil {
			// Log full stdout + stderr if yt-dlp fails
			output := string(outByte)
			srv.logger.Error("yt-dlp failed", "error", err, "output", output)
			sendResultError(fmt.Errorf("%s error: %v, output: %s", ytDlpName, err, output))
			return
		}

		// Debug log command output
		srv.logger.Debug(string(outByte))

		// Get the actual file size
		fileInfo, err := os.Stat(filePath)
		if err == nil {
			fileSize = uptr.Int(int(fileInfo.Size()))
		}

		// Build response struct
		result := &ddownload.DownloadResult{
			YoutubeTitle: title,
			FilePath:     filePath,
			Filename:     fileName,
			FileExt:      fileExt,
			FileFullName: FileFullName,
			Filesize:     fileSize,
		}

		resultCh <- result

		// Log successful download
		srv.logger.Info("Download successful", "info", result)
	}()

	return resultCh, nil
}

// prepareDownloadOptions prepare download options with defaults and user overrides
func (srv *YtDlpService) prepareDownloadOptions(
	options *ddownload.DownloadOptions,
) (
	formatType dtypes.FormatType,
	videoFormat dtypes.VideoFormat,
	audioFormat dtypes.AudioFormat,
	downloadDir, fileName string,
) {
	// Set default values
	formatType = formatTypeDefault
	videoFormat = videoFormatDefault
	audioFormat = audioFormatDefault
	downloadDir = srv.downloadsDir
	fileName = ""

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

	// Override download directory if provided
	if options.DownloadsDir != nil {
		downloadDir = *options.DownloadsDir
	}

	return formatType, videoFormat, audioFormat, downloadDir, fileName
}

// buildDownloadArgs build yt-dlp arguments and get file extension and title
func (srv *YtDlpService) buildDownloadArgs(
	url string,
	formatType dtypes.FormatType,
	videoFormat dtypes.VideoFormat,
	audioFormat dtypes.AudioFormat,
) (args []string, fileExt string, info *dyoutubeinfo.YouTubeInfo, err error) {

	// Default audio quality (used when extracting audio)
	var (
		audioQualityM4A = audioQualityM4ADefault
		audioQualityMP3 = audioQualityMP3Default
	)

	// prevent downloading the entire playlist, only fetch single video
	args = append(args, "--no-playlist")

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
			fileExt = info.Formats[0].FileExt
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
