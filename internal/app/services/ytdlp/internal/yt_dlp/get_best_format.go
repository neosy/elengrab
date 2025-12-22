package ytdlp

import (
	"context"
	"fmt"
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/yt_dlp/dto"
	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (y *YTDlp) GetBestFormat(
	ctx context.Context,
	url string,
	format string,
) (*dyoutubeinfo.YouTubeInfo, error) {
	info, err := y.getBestFormat(ctx, url, format)
	if err != nil {
		return nil, err
	}

	return y.mappers.VideoInfoToDomain(info), nil
}

func (y *YTDlp) getBestFormat(
	ctx context.Context,
	url string,
	format string,
) (*idto.YouTubeInfo, error) {
	err := y.updateCache(ctx, url)
	if err != nil {
		return nil, err
	}

	out, err := y.execCommandContext(
		ctx, y.ytDlpPath,
		"--no-warnings", "--quiet",
		"-f", format,
		"--load-info-json",
		y.formatCache.cacheFilePath(url),
		"--get-format",
	)
	if err != nil {
		y.formatCache.deleteByURL(url)
		return nil, err
	}

	outStr := strings.TrimSpace(string(out))
	if outStr == "" {
		y.formatCache.deleteByURL(url)
		return nil, fmt.Errorf("best format not found")
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
