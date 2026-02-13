package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/downloader/helper"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/utils"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (c *Core) GetBestFormat(
	ctx context.Context,
	url string,
	format string,
	useCookies bool,
) (*dmedia.MediaInfo, error) {
	info, err := c.getBestFormat(ctx, url, format, useCookies)
	if err != nil {
		return nil, err
	}

	return c.mappers.MapMediaInfoToDomain(info), nil
}

func (c *Core) getBestFormat(
	ctx context.Context,
	url string,
	format string,
	useCookies bool,
) (*idto.MediaInfo, error) {
	err := c.updateFormatCache(ctx, url, useCookies)
	if err != nil {
		return nil, err
	}

	// Prepare command arguments
	var args []string

	args = append(args, "--no-warnings", "--quiet")

	// Add YouTube cookies if allowed in service options
	if useCookies {
		args = helper.AddYouTubeCookiesToArgs(c.logger, args, c.serviceOptions)
	}

	args = append(args, "-f", format)
	args = append(args, "--load-info-json", c.formatCache.CacheFilePath(url))
	args = append(args, "--get-format")

	// Execute the command to get the best format ID
	out, err := utils.ExecCommandContext(ctx, c.ytDlpPath, args...)
	if err != nil {
		c.formatCache.DeleteByURL(url)
		return nil, err
	}

	outStr := strings.TrimSpace(string(out))
	if outStr == "" {
		c.formatCache.DeleteByURL(url)
		return nil, fmt.Errorf("best format not found")
	}

	formats := strings.SplitN(outStr, "+", 2)
	var bestFormatIds = make([]string, 0, len(formats))
	for _, f := range formats {
		parts := strings.SplitN(f, " - ", 2)
		bestFormatIds = append(bestFormatIds, parts[0])
	}

	info, err := c.getFormats(ctx, url, useCookies)
	if err != nil {
		return nil, err
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
