package ytdlp

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/google/uuid"
	iutils "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/utils"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/ytdlp/helper"
	iconstants "github.com/neosy/elengrab/internal/constants"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	"github.com/neosy/elengrab/pkg/nfile"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

const ()

func (y *YTDlp) Download(
	ctx context.Context,
	url string,
	concurrentFragments uint8,
	options *dservices.DownloadOptions,
	downloadResultCh chan<- *ddownload.DownloadResult,
) {
	type dlResult = ddownload.DownloadResult

	sendError := func(err error) {
		downloadResultCh <- &dlResult{Error: err}
	}

	sendData := func(data *dlResult) {
		downloadResultCh <- data
	}

	// Prepare download options with defaults and user overrides
	dlOptions, dlDir, fileName, includeTitleInFilename := helper.PrepareDownloadOptions(y.downloadsDir, concurrentFragments, options)

	// Ensure download directory exists
	if err := iutils.CheckDir(dlDir); err != nil {
		sendError(fmt.Errorf("failed to check directory: %w", err))
		return
	}

	title, err := y.getTitleFast(url)
	if err == nil && title != "" {
		sendData(&dlResult{
			YoutubeTitle: title,
		})
	}

	// Build yt-dlp arguments and get file extension and title
	args, fileExt, info, mediaInfo, err := helper.BuildDownloadArgs(
		ctx,
		url,
		dlOptions,
		y.GetBestFormat,
	)
	if err != nil {
		sendError(fmt.Errorf("failed to build download arguments: %w", err))
		return
	}

	// If title is empty, fetch it manually
	title = info.Title
	if title == "" {
		var err error
		title, err = y.GetTitle(ctx, url)
		if err != nil {
			sendError(fmt.Errorf("failed to get title: %w", err))
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
	filePath := path.Join(dlDir, FileFullName)

	var fileSize *int
	if len(info.Formats) == 1 {
		fileSize = info.Formats[0].Filesize
	}

	sendData(
		&dlResult{
			YoutubeTitle: title,
			FilePath:     filePath,
			Filename:     fileName,
			FileExt:      fileExt,
			FileFullName: FileFullName,
			Filesize:     fileSize,
		},
	)

	// Add output file path to yt-dlp arguments
	args = append(args, "-o", filePath)

	// Add the video URL to arguments
	args = append(args, url)

	// Execute yt-dlp command
	// Create command without CommandContext.
	// We manage cancellation manually to properly kill the whole process group.
	cmd := exec.Command(y.ytDlpPath, args...)

	// Running yt-dlp in a separate temporary directory
	baseTmpDir := filepath.Join(dlDir, ".yt-dlp")
	workDir, cleanup, err := iutils.CreateTempDir(baseTmpDir, "job-*")
	if err != nil {
		sendError(
			fmt.Errorf("%s failed to create tmp dir: %w", y.ytDlpName, err),
		)
		return
	}
	defer cleanup()

	cmd.Dir = workDir

	// Start the process in a new process group
	// so we can kill yt-dlp and all its children (ffmpeg).
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Capture stdout and stderr
	// Prepare stdout pipe
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		sendError(
			fmt.Errorf(
				"%s failed to create stdout pipe: %w",
				y.ytDlpName,
				err,
			),
		)
		return
	}

	// Prepare stderr pipe
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		sendError(
			fmt.Errorf(
				"%s failed to create stderr pipe: %w",
				y.ytDlpName,
				err,
			),
		)
		return
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		sendError(
			fmt.Errorf(
				"%s failed to start process: %w",
				y.ytDlpName,
				err,
			),
		)
		return
	}

	// Kill the entire process group on context cancellation
	go func() {
		<-ctx.Done()

		if cmd.Process == nil {
			return
		}

		pgid := -cmd.Process.Pid

		// Try graceful shutdown first
		_ = syscall.Kill(pgid, syscall.SIGTERM)

		// Wait a bit for yt-dlp / ffmpeg to cleanup temp files
		time.Sleep(2 * time.Second)

		// Force kill if still running
		if cmd.Process != nil {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
	}()

	// Read combined stdout + stderr
	out, _ := io.ReadAll(io.MultiReader(stdoutPipe, stderrPipe))

	// Wait for the process to exit
	err = cmd.Wait()
	if err != nil {
		// Context cancellation has priority
		if ctx.Err() != nil {
			sendError(
				fmt.Errorf("%s canceled: %w", y.ytDlpName, ctx.Err()),
			)
			return
		}

		// Process exited with an error
		sendError(
			fmt.Errorf("%s failed: %w, output: %s", y.ytDlpName, err, string(out)),
		)
		return
	}

	// Debug log command output
	y.logger.Debug(
		"Download completed",
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
	result := &dlResult{
		YoutubeTitle: title,
		FilePath:     filePath,
		Filename:     fileName,
		FileExt:      fileExt,
		FileFullName: FileFullName,
		Filesize:     fileSize,
		PartialHash:  partialHash,
		MediaInfo:    mediaInfo,
	}

	// Build response struct
	sendData(result)

	// Log successful download
	y.logger.Info("Download successful", "info", result)
}
