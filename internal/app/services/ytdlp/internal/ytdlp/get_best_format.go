package ytdlp

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/ytdlp/dto"
	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (y *YTDlp) GetBestFormat(ctx context.Context, url string, format string) (*dyoutubeinfo.YouTubeInfo, error) {
	info, err := y.getBestFormat(ctx, url, format)
	if err != nil {
		return nil, err
	}

	return y.mappers.VideoInfoToDomain(info), nil
}

func (y *YTDlp) getBestFormat(ctx context.Context, url string, format string) (*idto.YouTubeInfo, error) {
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

	info, err := y.getFormats(ctx, url)
	if err != nil {
		return nil, err
	}

	var bestFormat *idto.Format
	for _, f := range info.Formats {
		if f.FormatID == bestFormatId {
			bestFormat = &f
			break
		}
	}

	if bestFormat == nil {
		err := fmt.Errorf("best format not found")
		return nil, err
	}

	var bestInfo *idto.YouTubeInfo = uptr.Any(*info)
	bestInfo.Formats = []idto.Format{*bestFormat}

	return bestInfo, nil
}
