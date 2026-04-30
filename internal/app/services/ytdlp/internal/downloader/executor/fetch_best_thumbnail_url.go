package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
)

func (e *Executor) ExtractBestThumbnailURL(ctx context.Context, mediaURL string, opts ...Option) (string, error) {
	var (
		opt = newDefaultOptions()
		cmd *exec.Cmd
	)

	for _, o := range opts {
		o(opt)
	}

	if opt.ensureCache {
		err := e.EnsureFormatCache(ctx, mediaURL, opt.useCookies)
		if err != nil {
			return "", fmt.Errorf("failed to ensure format cache: %w", err)
		}
	}

	// Prepare command arguments
	var args []string

	args = append(args, "--no-playlist", "--no-warnings", "--print", "%(thumbnails.-1.url)s")

	// Add YouTube cookies if allowed in service options
	if opt.useCookies {
		args = addYouTubeCookiesToArgs(e.logger, args, e.serviceOptions)
	}

	args = append(args, "--load-info-json", e.formatCache.CacheFilePath(mediaURL))
	cmd = exec.CommandContext(ctx, e.ytDlpPath, args...)

	// Buffers to capture stdout and stderr
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// The context was canceled
		if ctx.Err() != nil {
			return "", fmt.Errorf("process canceled: %w", ctx.Err())
		}
		errOut := fmt.Errorf("%s failed: %w, stderr: %s", consts.YtDlpName, err, stderr.String())
		return "", errOut
	}

	thumbnailURL := strings.TrimSpace(out.String())
	if thumbnailURL == "" {
		err := fmt.Errorf("thumbnail not found")
		return "", err
	}

	return thumbnailURL, nil
}
