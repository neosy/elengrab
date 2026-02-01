package downloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/utils"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/pkg/syncx"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (d *Downloader) runYtDlp(
	ctx context.Context,
	url string,
	meta *idto.DownloadMeta,
	onProgressUpdate func(ddownload.DownloadProgress),
) ([]byte, error) {
	var (
		done      = syncx.NewDoneSignal()
		isTimeOut atomic.Bool
	)
	defer done.Close()

	dlDir := filepath.Dir(meta.FilePath)

	// Cache directory
	cacheDir := filepath.Join(dlDir, ytDlpCacheDir)

	// Running yt-dlp in a separate temporary directory
	baseTmpDir := filepath.Join(dlDir, ytDlpTempDir)

	workDir, cleanup, err := utils.CreateTempDir(baseTmpDir, "job-*")
	if err != nil {
		return nil, fmt.Errorf("%s failed to create tmp dir: %w", d.ytDlpName, err)
	}
	defer func() {
		err := cleanup()
		if err != nil {
			d.logger.Debug("Failed clear temp dir", "error", err)
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
		args = append(args, "--load-info-json", d.formatCache.CacheFilePath(url))
	} else {
		// Add extractor args
		args = append(args, "--extractor-args", *meta.Options.ExtractorArgs)

		// Do not use cache for custom extractor args
		args = append(args, meta.URL)
	}

	// Add output file path to yt-dlp arguments
	args = append(args, "-o", tmpFilePath)

	// Execute yt-dlp command
	// Create command without CommandContext.
	// We manage cancellation manually to properly kill the whole process group.
	cmd := exec.Command(d.ytDlpPath, args...)

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
			d.ytDlpName,
			err,
		)
	}

	// Prepare stderr pipe
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf(
			"%s failed to create stderr pipe: %w",
			d.ytDlpName,
			err,
		)
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf(
			"%s failed to start process: %w",
			d.ytDlpName,
			err,
		)
	}

	// Kill the entire process group on context cancellation
	go func() {
		timer := time.NewTimer(ytDlpTimeout)
		defer timer.Stop()

		select {
		case <-done.Done():
			return
		case <-ctx.Done():
			d.logger.Debug(
				fmt.Sprintf("context canceled, killing process %s", d.ytDlpName),
				"url", url,
			)
		case <-timer.C:
			d.logger.Warn(
				fmt.Sprintf("%s timeout reached, killing process", d.ytDlpName),
				"url", url,
			)
			isTimeOut.Store(true)
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

	var wg sync.WaitGroup

	// stdout reader (progress)
	var outBuf bytes.Buffer
	wg.Go(func() {
		var fileSize *int64
		if meta.FileSize != nil {
			fileSize = uptr.Any(int64(*meta.FileSize))
		}
		// Watch progress output
		d.watchProgress(fileSize, stdoutPipe, &outBuf, onProgressUpdate)
	})

	// stderr reader
	var errBuf bytes.Buffer
	wg.Go(func() {
		_, _ = io.Copy(&errBuf, stderrPipe)
	})

	// Wait for the process to exit
	err = cmd.Wait()
	done.Close()
	wg.Wait()
	if err != nil {
		// Context cancellation has priority
		if ctx.Err() != nil {
			return nil, fmt.Errorf("%s canceled: %w", d.ytDlpName, ctx.Err())
		}

		// Check if timeout occurred
		if isTimeOut.Load() {
			return nil, fmt.Errorf("%s timeout reached", d.ytDlpName)
		}

		// Process exited with an error
		return nil, fmt.Errorf("%s failed: %w, output: %s", d.ytDlpName, err, errBuf.String())
	}

	// Check if timeout occurred
	if isTimeOut.Load() {
		return nil, fmt.Errorf("%s timeout reached", d.ytDlpName)
	}

	// Move final file from temp directory to target path
	// to avoid ffmpeg creating temporary files in the download directory
	err = os.Rename(tmpFilePath, meta.FilePath)
	if err != nil {
		return nil, fmt.Errorf("%s failed to move file from temp dir to target path: %w", d.ytDlpName, err)
	}

	return outBuf.Bytes(), nil
}
