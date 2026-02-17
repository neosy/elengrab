package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/utils"
)

func (e *Executor) fetchFormatsJSON(
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
		args = helper.AddYouTubeCookiesToArgs(e.logger, args, e.serviceOptions)
	}

	args = append(args, url)

	// Execute the command to get video formats in JSON format
	out, err := utils.ExecCommandContext(ctx, e.ytDlpPath, args...)
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

func (e *Executor) fetchAndCacheFormatsJSON(
	ctx context.Context,
	url string,
	useCookies bool,
) ([]byte, error) {
	dataJSON, err := e.fetchFormatsJSON(ctx, url, useCookies)
	if err != nil {
		e.formatCache.DeleteByURL(url)
		e.logger.Debug("Failed to fetch formats JSON", "url", url, "error", err)
		return nil, err
	}
	e.logger.Debug("Fetched formats JSON", "url", url)

	err = e.formatCache.WriteByURL(url, dataJSON)
	if err != nil {
		return nil, err
	}

	e.logger.Debug("Updated formats JSON cache", "url", url)

	return dataJSON, nil
}

func (e *Executor) getFormatsJSON(
	ctx context.Context,
	url string,
	useCookies bool,
) ([]byte, error) {
	dataJSON, err := e.formatCache.LoadByURL(url)
	if err != nil {
		return nil, err
	}

	if dataJSON != nil {
		return dataJSON, nil
	}

	dataJSON, err = e.fetchAndCacheFormatsJSON(ctx, url, useCookies)
	if err != nil {
		return nil, err
	}

	return dataJSON, nil
}

func (e *Executor) GetFormats(
	ctx context.Context,
	url string,
	useCookies bool,
) (*idto.MediaInfo, error) {
	dataJSON, err := e.getFormatsJSON(ctx, url, useCookies)
	if err != nil {
		return nil, fmt.Errorf("get formats json error: %w", err)
	}

	var info = &idto.MediaInfo{}
	err = json.NewDecoder(bytes.NewReader(dataJSON)).Decode(info)
	if err != nil {
		e.formatCache.DeleteByURL(url)
		return nil, fmt.Errorf("json decode error: %w", err)
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
