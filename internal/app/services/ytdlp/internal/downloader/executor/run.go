package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/utils"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	"github.com/neosy/elengrab/internal/pkg/syncx"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (e *Executor) RunYtDlp(
	ctx context.Context,
	url string,
	meta *idto.DownloadMeta,
	onProgressUpdate func(dservices.DownloaderProgress),
) ([]byte, error) {
	var (
		done      = syncx.NewDoneSignal()
		isTimeOut atomic.Bool
	)
	defer done.Close()

	// Cache directory
	cacheDir := filepath.Join(e.storage.BasePath(), consts.YtDlpCacheDir)

	// Running yt-dlp in a separate temporary directory
	baseTmpDir := filepath.Join(e.storage.BasePath(), consts.YtDlpTempDir)

	workDir, cleanup, err := utils.CreateTempDir(baseTmpDir, "job-*")
	if err != nil {
		return nil, fmt.Errorf("%s failed to create tmp dir: %w", consts.YtDlpName, err)
	}
	defer func() {
		err := cleanup()
		if err != nil {
			e.logger.Debug("Failed clear temp dir", "error", err)
		}
	}()

	// Copy args to avoid modifying the original slice
	args := meta.Options.Args[0:len(meta.Options.Args):len(meta.Options.Args)]

	// Add concurrent fragments argument
	args = append(args, "--concurrent-fragments", strconv.Itoa(int(meta.Options.ConcurrentFragments)))

	// Add progress output to yt-dlp arguments
	args = append(
		args,
		"--newline",
		"--progress-template",
		"%(progress.downloaded_bytes)s|%(progress.total_bytes)s|%(progress.eta)s|%(progress.speed)s",
	)

	// Build full path to the output file inside the temp work directory
	tmpFilePath := filepath.Join(workDir, meta.FileFullName)

	// If no extractor args provided, use cache and isolated paths
	if meta.Options.ExtractorArgs == nil {
		// Add cache directory to yt-dlp arguments
		args = append(args, "--cache-dir", cacheDir)

		// Force yt-dlp to store all temporary and intermediate files in the isolated work directory
		args = append(args, "--paths", fmt.Sprintf("temp:%s", workDir))

		// Add load info json to arguments
		args = append(args, "--load-info-json", e.formatCache.CacheFilePath(url))
	} else {
		// Add extractor args
		args = append(args, "--extractor-args", *meta.Options.ExtractorArgs)

		// Do not use cache for custom extractor args
		args = append(args, meta.URL)
	}

	// Add YouTube cookies if allowed in service options
	if meta.Options.RequiresYouTubeCookies {
		args = addYouTubeCookiesToArgs(e.logger, args, e.serviceOptions)
	}

	// Add output file path to yt-dlp arguments
	args = append(args, "-o", tmpFilePath)

	// Execute yt-dlp command
	// Create command without CommandContext.
	// We manage cancellation manually to properly kill the whole process group.
	cmd := exec.Command(e.ytDlpPath, args...)

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
			consts.YtDlpName,
			err,
		)
	}

	// Prepare stderr pipe
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf(
			"%s failed to create stderr pipe: %w",
			consts.YtDlpName,
			err,
		)
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf(
			"%s failed to start process: %w",
			consts.YtDlpName,
			err,
		)
	}

	var wg sync.WaitGroup

	// Kill the entire process group on context cancellation
	wg.Go(func() {
		timer := time.NewTimer(consts.YtDlpTimeout)
		defer timer.Stop()

		select {
		case <-done.Done():
			e.logger.Debug(
				"Done signal received, cleaning up processes",
				"name", consts.YtDlpName,
				"url", url,
			)
		case <-ctx.Done():
			e.logger.Debug(
				"Context canceled, killing process",
				"name", consts.YtDlpName,
				"url", url,
			)
		case <-timer.C:
			e.logger.Warn(
				"Timeout reached, killing process",
				"name", consts.YtDlpName,
				"url", url,
			)
			isTimeOut.Store(true)
		}

		if cmd.Process == nil {
			e.logger.Debug(
				"Process already exited, nothing to kill",
				"name", consts.YtDlpName,
				"url", url,
			)
			return
		}

		pgid := cmd.Process.Pid

		// Try graceful shutdown first
		tryGracefulKill(-pgid)

		// Wait a bit for yt-dlp / ffmpeg to cleanup temp files
		select {
		case <-time.After(500 * time.Millisecond):
			// wait
		case <-ctx.Done():
			// force kill
		}

		// Force kill if still running
		forceKill(cmd)
	})

	// stdout reader (progress)
	var outBuf bytes.Buffer
	wg.Go(func() {
		var fileSize *int64
		if meta.FileSize != nil {
			fileSize = uptr.Any(int64(*meta.FileSize))
		}
		// Watch progress output
		e.watchProgress(fileSize, stdoutPipe, &outBuf, onProgressUpdate)
	})

	// stderr reader
	var errBuf bytes.Buffer
	wg.Go(func() {
		_, _ = io.Copy(&errBuf, stderrPipe)
	})

	// Wait for the process (yt-dlp) to exit
	err = cmd.Wait()

	// Complete the goroutines
	done.Close()

	// Waiting for goroutines to complete
	wg.Wait()

	if err != nil {
		// Context cancellation has priority
		if ctx.Err() != nil {
			return nil, fmt.Errorf("%s canceled: %w", consts.YtDlpName, ctx.Err())
		}

		// Check if timeout occurred
		if isTimeOut.Load() {
			return nil, fmt.Errorf("%s timeout reached", consts.YtDlpName)
		}

		// Process exited with an error
		return nil, fmt.Errorf("%s failed: %w, output: %s", consts.YtDlpName, err, errBuf.String())
	}

	// Check if timeout occurred
	if isTimeOut.Load() {
		return nil, fmt.Errorf("%s timeout reached", consts.YtDlpName)
	}

	// Move final file from temp directory to target path
	// to avoid ffmpeg creating temporary files in the download directory
	err = e.storage.Move(tmpFilePath, meta.FileFullName)
	if err != nil {
		return nil, fmt.Errorf("%s failed to move file from temp dir to target path: %w", consts.YtDlpName, err)
	}

	return outBuf.Bytes(), nil
}
