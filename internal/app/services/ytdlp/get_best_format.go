package ytdlpsrv

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
)

// GetBestFormat retrieves and parses video best format for the given URL.
func (srv *YtDlpService) GetBestFormat(ctx context.Context, url string) (*dyoutubeinfo.YouTubeInfo, error) {
	return srv.getBestFormat(ctx, url, "b")
}

func (srv *YtDlpService) getBestFormat(ctx context.Context, url string, format string) (*dyoutubeinfo.YouTubeInfo, error) {
	cmd := exec.CommandContext(ctx, srv.cmdPath, "--no-playlist", "--no-warnings", "-f", format, "--get-format", url)

	// Buffers to capture stdout and stderr
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// The context was canceled
		if ctx.Err() != nil {
			return nil, fmt.Errorf("process canceled: %w", ctx.Err())
		}
		errOut := fmt.Errorf("%s failed: %v, stderr: %s", ytDlpName, err, stderr.String())
		return nil, errOut
	}

	outStr := strings.TrimSpace(out.String())
	if outStr == "" {
		err := fmt.Errorf("best format not found")
		return nil, err
	}

	parts := strings.SplitN(outStr, " - ", 2)
	bestFormatId := parts[0]

	info, err := srv.GetFormats(ctx, url)
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

	srv.logger.Debug(
		"Get best format",
		"url", url,
		"format", bestInfo,
	)

	return &bestInfo, nil
}
