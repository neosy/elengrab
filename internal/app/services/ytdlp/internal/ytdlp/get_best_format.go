package ytdlp

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
)

func (y *YTDlp) GetBestFormat(ctx context.Context, url string, format string) (*dyoutubeinfo.YouTubeInfo, error) {
	cmd := exec.CommandContext(ctx, y.ytDlpPath, "--no-playlist", "--no-warnings", "-f", format, "--get-format", url)

	// Buffers to capture stdout and stderr
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// The context was canceled
		if ctx.Err() != nil {
			return nil, fmt.Errorf("process canceled: %w", ctx.Err())
		}
		errOut := fmt.Errorf("%s failed: %v, stderr: %s", y.ytDlpName, err, stderr.String())
		return nil, errOut
	}

	outStr := strings.TrimSpace(out.String())
	if outStr == "" {
		err := fmt.Errorf("best format not found")
		return nil, err
	}

	parts := strings.SplitN(outStr, " - ", 2)
	bestFormatId := parts[0]

	info, err := y.GetFormats(ctx, url)
	if err != nil {
		return nil, err
	}

	var bestFormat *dyoutubeinfo.Format
	for _, f := range info.Formats {
		if f.FormatId == bestFormatId {
			bestFormat = &f
			break
		}
	}

	if bestFormat == nil {
		err := fmt.Errorf("best format not found")
		return nil, err
	}

	bestInfo := *info
	bestInfo.Formats = []dyoutubeinfo.Format{*bestFormat}

	return &bestInfo, nil
}
