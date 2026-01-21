package ytdlp

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/dto"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/yt_dlp/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/yt_dlp/helper"
	"github.com/neosy/elengrab/internal/app/utils"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/neosy/elengrab/pkg/nfile"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (y *YTDlp) Download(
	ctx context.Context,
	url string,
	concurrentFragments uint8,
	options *dservices.DownloadOptions,
	downloadResultCh chan<- *ddownload.DownloadResult,
) {
	var wg sync.WaitGroup
	doneCh := make(chan struct{})
	defer func() {
		close(doneCh)
		wg.Wait()
	}()

	var (
		sendError = func(data *ddownload.DownloadResult, err error) {
			var result = &ddownload.DownloadResult{}
			if data != nil {
				*result = *data
			}
			result.Error = err

			select {
			case <-doneCh:
			case downloadResultCh <- result:
			case <-ctx.Done():
			}
		}

		sendData = func(data *ddownload.DownloadResult) {
			select {
			case <-doneCh:
			case downloadResultCh <- data:
			case <-ctx.Done():
			}
		}
	)

	url = strings.TrimSpace(url)

	// Prepare download options with defaults and user overrides
	dlOptions, dlDir, fileName, includeTitleInFilename :=
		helper.PrepareDownloadOptions(y.downloadsDir, concurrentFragments, options)

	// Ensure download directory exists
	if err := nfile.CheckDir(dlDir); err != nil {
		sendError(nil, fmt.Errorf("failed to check directory: %w", err))
		return
	}

	title, err := y.getTitleFast(url)
	if err == nil && title != "" {
		sendData(&ddownload.DownloadResult{
			YoutubeTitle: title,
		})
	}

	meta, err := y.prepareMetadata(ctx, url, dlDir, fileName, includeTitleInFilename, dlOptions)
	if err != nil {
		sendError(&ddownload.DownloadResult{YoutubeTitle: title}, err)
		return
	}

	sendData(meta.InitialResult())

	// Start asynchronous fetching of the channel avatar.
	// Returns a channel from which the avatar can be read once the goroutine completes.
	channelAvatarCh := y.fetchChannelAvatarAsync(&wg, *meta, sendData)

	// Run yt-dlp for the given URL and metadata.
	// Capture output and error; send error if execution fails.
	out, err := y.runYtDlp(ctx, url, meta)
	if err != nil {
		sendError(meta.InitialResult(), err)
		return
	}

	// Debug log command output
	y.logger.Debug(
		"Download completed",
		"url", url,
		"out", string(out),
	)

	// Get the actual file size
	fileInfo, err := os.Stat(meta.FilePath)
	if err == nil {
		meta.FileSize = uptr.Int(int(fileInfo.Size()))
	}

	var partialHash *string
	{
		h, err := utils.HashPartialMedia(meta.FilePath)
		if err == nil && h != "" {
			partialHash = &h
		}
	}

	// Waiting for background processes to complete
	wg.Wait()

	// Build response struct
	result := meta.InitialResult()
	{
		result.PartialHash = partialHash
		result.ChannelAvatar = <-channelAvatarCh
	}

	// Build response struct
	sendData(result)

	// Log successful download
	// Info Download completed
	y.logger.Info(
		"Download completed",
		"title", meta.Title,
		"url", url,
		"mediaInfo", result.MediaInfo,
	)
	y.logger.Debug("Download success", "meta", meta)
}

func (y *YTDlp) prepareMetadata(
	ctx context.Context,
	url, dlDir, fileName string,
	includeTitleInFilename bool,
	options dto.DLOptions,
) (*idto.DownloadMeta, error) {
	// Build yt-dlp arguments and get file extension and title
	args, fileExt, info, mediaInfo, err := helper.BuildDownloadArgs(
		ctx,
		url,
		options,
		y.getBestFormat,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build download arguments: %w", err)
	}

	// If title is empty, fetch it manually
	title := info.Title
	if title == "" {
		var err error
		title, err = y.GetTitle(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("failed to get title: %w", err)
		}
	}

	// Generate a unique filename if none provided
	if fileName == "" {
		fileName = uuid.New().String()
	}

	if includeTitleInFilename {
		fileName = fmt.Sprintf("%s_%s", nfasthttp.SanitizeFileName(title), fileName)
	}

	fileFullName := fmt.Sprintf("%s.%s", fileName, fileExt)
	filePath := filepath.Join(dlDir, fileFullName)

	var (
		fileSize *int
		tmpSize  int
	)
	for _, f := range info.Formats {
		if f.Filesize != nil {
			tmpSize += *f.Filesize
		}
	}
	if tmpSize != 0 {
		fileSize = &tmpSize
	}

	var channelID *string
	if info.ChannelID != "" {
		channelID = &info.ChannelID
	}

	return &idto.DownloadMeta{
		Title:        title,
		FileName:     fileName,
		FileExt:      fileExt,
		FileFullName: fileFullName,
		FilePath:     filePath,
		FileSize:     fileSize,
		ChannelID:    channelID,
		ChannelURL:   info.ChannelUrl,
		MediaInfo:    mediaInfo,
		Args:         args,
	}, nil
}

func (y *YTDlp) fetchChannelAvatarAsync(
	wg *sync.WaitGroup,
	meta idto.DownloadMeta,
	send func(*ddownload.DownloadResult),
) <-chan *ddownload.DownloadResultChannelAvatar {
	channelAvatar := make(chan *ddownload.DownloadResultChannelAvatar, 1)

	if meta.ChannelID == nil {
		close(channelAvatar)
		return channelAvatar
	}

	wg.Go(func() {
		defer close(channelAvatar)

		avatarSources, err := y.getChannelAvatar(meta.ChannelURL)
		if err != nil {
			y.logger.Debug("Failed to get channel avatar", "channelURL", meta.ChannelURL, "error", err)
			return
		}
		if len(avatarSources) == 0 {
			y.logger.Debug("Avatar not found", "channelURL", meta.ChannelURL)
			return
		}

		var avatar *ddownload.DownloadResultChannelAvatar
		src := avatarSources[0]
		if len(src.Raw) == 0 {
			y.logger.Debug("Avatar image not found", "channelURL", meta.ChannelURL)
			return
		}

		y.logger.Info(
			"YouTube channel avatar fetched successfully",
			"channelURL", meta.ChannelURL,
		)

		avatar = &ddownload.DownloadResultChannelAvatar{
			ImageURL:    src.URL,
			ImageRAW:    src.Raw,
			ImageFormat: src.Format,
		}
		channelAvatar <- avatar

		result := meta.InitialResult()
		result.ChannelAvatar = avatar

		send(result)
	})

	return channelAvatar
}

func (y *YTDlp) runYtDlp(
	ctx context.Context,
	url string,
	meta *idto.DownloadMeta,
) ([]byte, error) {
	var doneCh = make(chan struct{})
	defer close(doneCh)

	dlDir := filepath.Dir(meta.FilePath)

	// Cache directory
	cacheDir := filepath.Join(dlDir, ytDlpCacheDir)

	// Running yt-dlp in a separate temporary directory
	baseTmpDir := filepath.Join(dlDir, ytDlpTempDir)

	workDir, cleanup, err := helper.CreateTempDir(baseTmpDir, "job-*")
	if err != nil {
		return nil, fmt.Errorf("%s failed to create tmp dir: %w", y.ytDlpName, err)
	}
	defer func() {
		err := cleanup()
		if err != nil {
			y.logger.Debug("Failed clear temp dir", "error", err)
		}
	}()

	// Build full path to the output file inside the temp work directory
	tmpFilePath := filepath.Join(workDir, meta.FileFullName)

	// Add cache directory to yt-dlp arguments
	meta.Args = append(meta.Args, "--cache-dir", cacheDir)

	// Force yt-dlp to store all temporary and intermediate files in the isolated work directory
	meta.Args = append(meta.Args, "--paths", fmt.Sprintf("temp:%s", workDir))

	// Add load info json to arguments
	meta.Args = append(meta.Args, "--load-info-json", y.formatCache.cacheFilePath(url))

	// Add output file path to yt-dlp arguments
	meta.Args = append(meta.Args, "-o", tmpFilePath)

	// Execute yt-dlp command
	// Create command without CommandContext.
	// We manage cancellation manually to properly kill the whole process group.
	cmd := exec.Command(y.ytDlpPath, meta.Args...)

	cmd.Dir = workDir

	// Start the process in a new process group
	// so we can kill yt-dlp and all its children (ffmpeg).
	cmd.SysProcAttr = newSysProcAttr()

	// Capture stdout and stderr
	// Prepare stdout pipe
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf(
			"%s failed to create stdout pipe: %w",
			y.ytDlpName,
			err,
		)
	}

	// Prepare stderr pipe
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf(
			"%s failed to create stderr pipe: %w",
			y.ytDlpName,
			err,
		)
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf(
			"%s failed to start process: %w",
			y.ytDlpName,
			err,
		)
	}

	// Kill the entire process group on context cancellation
	go func() {
		select {
		case <-doneCh:
			return
		case <-ctx.Done():
		}

		if cmd.Process == nil {
			return
		}

		pgid := -cmd.Process.Pid

		// Try graceful shutdown first
		tryGracefulKill(pgid)

		// Wait a bit for yt-dlp / ffmpeg to cleanup temp files
		time.Sleep(2 * time.Second)

		// Force kill if still running
		forceKill(cmd)
	}()

	// Read combined stdout + stderr
	out, _ := io.ReadAll(io.MultiReader(stdoutPipe, stderrPipe))

	// Wait for the process to exit
	err = cmd.Wait()
	if err != nil {
		// Context cancellation has priority
		if ctx.Err() != nil {
			return nil, fmt.Errorf("%s canceled: %w", y.ytDlpName, ctx.Err())
		}

		// Deleting the cache, because Youtube could have changed the format
		y.formatCache.deleteByURL(url)

		// Process exited with an error
		return nil, fmt.Errorf("%s failed: %w, output: %s", y.ytDlpName, err, string(out))
	}

	// Move final file from temp directory to target path
	// to avoid ffmpeg creating temporary files in the download directory
	err = os.Rename(tmpFilePath, meta.FilePath)
	if err != nil {
		return nil, fmt.Errorf("%s failed to move file from temp dir to target path: %w", y.ytDlpName, err)
	}

	return out, nil
}
