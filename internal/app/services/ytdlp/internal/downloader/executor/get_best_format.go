package executor

import (
	"context"
	"fmt"
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/utils"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (e *Executor) GetBestFormat(
	ctx context.Context,
	url string,
	format string,
	useCookies bool,
) (*idto.MediaInfo, error) {
	err := e.UpdateFormatCache(ctx, url, useCookies)
	if err != nil {
		return nil, err
	}

	// Prepare command arguments
	var args []string

	args = append(args, "--no-warnings", "--quiet")

	// Add YouTube cookies if allowed in service options
	if useCookies {
		args = helper.AddYouTubeCookiesToArgs(e.logger, args, e.serviceOptions)
	}

	args = append(args, "-f", format)
	args = append(args, "--load-info-json", e.formatCache.CacheFilePath(url))
	args = append(args, "--get-format")

	// Execute the command to get the best format ID
	out, err := utils.ExecCommandContext(ctx, e.ytDlpPath, args...)
	if err != nil {
		e.formatCache.DeleteByURL(url)
		return nil, fmt.Errorf("exec command error: %w", err)
	}

	outStr := strings.TrimSpace(string(out))
	if outStr == "" {
		e.formatCache.DeleteByURL(url)
		return nil, fmt.Errorf("best format not found")
	}

	formats := strings.SplitN(outStr, "+", 2)
	var bestFormatIds = make([]string, 0, len(formats))
	for _, f := range formats {
		parts := strings.SplitN(f, " - ", 2)
		bestFormatIds = append(bestFormatIds, parts[0])
	}

	info, err := e.GetFormats(ctx, url, useCookies)
	if err != nil {
		return nil, fmt.Errorf("get formats error: %w", err)
	}

	var bestFormats []idto.MediaFormat
	for _, bf := range bestFormatIds {
		for _, f := range info.Formats {
			if bf == f.FormatID {
				bestFormats = append(bestFormats, f)
			}
		}
	}

	if len(bestFormats) == 0 {
		err := fmt.Errorf("best format not found")
		return nil, err
	}

	bestInfo := uptr.Copy(info)
	bestInfo.Formats = bestFormats

	return bestInfo, nil
}
