package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/downloader/helper"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/core/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/utils"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

// GetFormats retrieves and parses video formats for the given URL.
func (c *Core) GetFormats(ctx context.Context, url string, useCookies bool) (*dmedia.MediaInfo, error) {
	info, err := c.getFormats(ctx, url, useCookies)
	if err != nil {
		return nil, err
	}

	return c.mappers.MapMediaInfoToDomain(info), nil
}

func (c *Core) fetchFormatsJSON(
	ctx context.Context,
	url string,
	useCookies bool,
) ([]byte, error) {
	// Prepare command arguments
	var args []string

	args = append(args, "--no-playlist", "--no-warnings", "--quiet")
	args = append(args, "--dump-json")

	// Add YouTube cookies if allowed in service options
	if useCookies {
		args = helper.AddYouTubeCookiesToArgs(c.logger, args, c.serviceOptions)
	}

	args = append(args, url)

	// Execute the command to get video formats in JSON format
	out, err := utils.ExecCommandContext(ctx, c.ytDlpPath, args...)
	if err != nil {
		return nil, err
	}

	outStr := string(out)

	// Find JSON start
	start := strings.Index(outStr, "{")
	if start == -1 {
		err := errors.New("no JSON found in yt-dlp output")
		return nil, err
	}

	dataJSON := []byte(outStr[start:])

	var js json.RawMessage

	if err := json.Unmarshal(dataJSON, &js); err != nil {
		return nil, fmt.Errorf("invalid JSON output from yt-dlp: %w", err)
	}

	return dataJSON, nil
}

func (c *Core) fetchAndCacheFormatsJSON(
	ctx context.Context,
	url string,
	useCookies bool,
) ([]byte, error) {
	dataJSON, err := c.fetchFormatsJSON(ctx, url, useCookies)
	if err != nil {
		c.formatCache.DeleteByURL(url)
		return nil, err
	}

	err = c.formatCache.WriteByURL(url, dataJSON)
	if err != nil {
		return nil, err
	}

	return dataJSON, nil
}

func (c *Core) getFormatsJSON(
	ctx context.Context,
	url string,
	useCookies bool,
) ([]byte, error) {
	dataJSON, err := c.formatCache.LoadByURL(url)
	if err != nil {
		return nil, err
	}

	if dataJSON != nil {
		return dataJSON, nil
	}

	dataJSON, err = c.fetchAndCacheFormatsJSON(ctx, url, useCookies)
	if err != nil {
		return nil, err
	}

	return dataJSON, nil
}

func (c *Core) getFormats(
	ctx context.Context,
	url string,
	useCookies bool,
) (*idto.MediaInfo, error) {
	dataJSON, err := c.getFormatsJSON(ctx, url, useCookies)
	if err != nil {
		return nil, err
	}

	var info = &idto.MediaInfo{}
	err = json.NewDecoder(bytes.NewReader(dataJSON)).Decode(info)
	if err != nil {
		c.formatCache.DeleteByURL(url)
		return nil, fmt.Errorf("json decode error: %v", err)
	}
	for i, f := range info.Formats {
		if f.Tbr != 0 && info.Duration != 0 {
			size := int64(int64(f.Tbr*1000) / 8 * int64(info.Duration))
			info.Formats[i].FilesizeApprox = &size
		}
	}

	info.ChannelTitle = strings.TrimSpace(info.ChannelTitle)

	return info, nil
}
