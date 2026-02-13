package core

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/downloader/helper"
)

func (c *Core) GetTitle(ctx context.Context, url string, useCookies bool) (string, error) {
	var cmd *exec.Cmd

	// Prepare command arguments
	var args []string

	if ok, _ := c.formatCache.IsTTLValidByURL(url); ok {
		args = append(args, "--no-playlist", "--no-warnings", "-e")

		// Add YouTube cookies if allowed in service options
		if useCookies {
			args = helper.AddYouTubeCookiesToArgs(c.logger, args, c.serviceOptions)
		}

		args = append(args, "--load-info-json", c.formatCache.CacheFilePath(url))
		cmd = exec.CommandContext(ctx, c.ytDlpPath, args...)
	} else {
		args = append(args, "--no-playlist", "--no-warnings", "-e")

		// Add YouTube cookies if allowed in service options
		if useCookies {
			args = helper.AddYouTubeCookiesToArgs(c.logger, args, c.serviceOptions)
		}

		args = append(args, url)
		cmd = exec.CommandContext(ctx, c.ytDlpPath, args...)
	}

	// Buffers to capture stdout and stderr
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// The context was canceled
		if ctx.Err() != nil {
			return "", fmt.Errorf("process canceled: %w", ctx.Err())
		}
		errOut := fmt.Errorf("%s failed: %v, stderr: %s", c.ytDlpName, err, stderr.String())
		return "", errOut
	}

	title := strings.TrimSpace(out.String())
	if title == "" {
		err := fmt.Errorf("title not found")
		return "", err
	}

	return title, nil
}
