package ytdlp

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func (y *YTDlp) GetTitle(ctx context.Context, url string) (string, error) {
	cmd := exec.CommandContext(ctx, y.ytDlpPath, "--no-playlist", "--no-warnings", "-e", url)

	// Buffers to capture stdout and stderr
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// The context was canceled
		if ctx.Err() != nil {
			return "", fmt.Errorf("process canceled: %w", ctx.Err())
		}
		errOut := fmt.Errorf("%s failed: %v, stderr: %s", y.ytDlpName, err, stderr.String())
		return "", errOut
	}

	title := strings.TrimSpace(out.String())
	if title == "" {
		err := fmt.Errorf("title not found")
		return "", err
	}

	return title, nil
}

func (y *YTDlp) getTitleFast(url string) (string, error) {
	info, err := y.getInfoFast(url)
	if err != nil {
		y.logger.Debug(
			"Get title fast",
			"url", url,
			"error", err,
		)
		return "", err
	}

	y.logger.Debug(
		"Get title fast",
		"url", url,
		"title", info.Title,
	)

	return info.Title, nil
}
