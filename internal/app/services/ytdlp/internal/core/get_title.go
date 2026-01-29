package core

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func (c *Core) GetTitle(ctx context.Context, url string) (string, error) {
	var cmd *exec.Cmd

	if ok, _ := c.formatCache.IsTTLValidByURL(url); ok {
		cmd = exec.CommandContext(
			ctx, c.ytDlpPath,
			"--no-playlist", "--no-warnings",
			"-e",
			"--load-info-json", c.formatCache.CacheFilePath(url),
		)
	} else {
		cmd = exec.CommandContext(ctx, c.ytDlpPath, "--no-playlist", "--no-warnings", "-e", url)
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
