package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/utils"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

func (e *Executor) FetchInfo(
	ctx context.Context,
	url string,
	opts ...idto.ExecutorOption,
) (*idto.ExtractInfo, error) {
	dataJSON, err := e.loadInfoJSON(ctx, url, opts...)
	if err != nil {
		return nil, fmt.Errorf("get formats json error: %w", err)
	}

	var info = &idto.ExtractInfo{}
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

	if hostdetect.Instagram(url) && info.Description != "" {
		first, _, _ := strings.Cut(info.Description, "\n")
		if first != "" {
			info.Title = first
		}
	}

	return info, nil
}

func (e *Executor) loadInfoJSON(
	ctx context.Context,
	url string,
	opts ...idto.ExecutorOption,
) ([]byte, error) {
	startTime := time.Now()
	dataJSON, err := e.formatCache.LoadByURL(url)
	elapsed := time.Since(startTime)
	if err != nil {
		return nil, err
	}

	if dataJSON != nil {
		e.logger.Debug(
			"Loaded formats JSON from cache",
			"url", url,
			"elapsed", uformat.DurationFormat(elapsed),
		)
		return dataJSON, nil
	}

	dataJSON, err = e.fetchAndCacheInfoJSON(ctx, url, opts...)
	if err != nil {
		return nil, err
	}

	return dataJSON, nil
}

func (e *Executor) fetchAndCacheInfoJSON(
	ctx context.Context,
	url string,
	opts ...idto.ExecutorOption,
) ([]byte, error) {
	startTime := time.Now()
	dataJSON, err := e.fetchInfoJSON(ctx, url, opts...)
	elapsed := time.Since(startTime)
	if err != nil {
		e.formatCache.DeleteByURL(url)
		e.logger.Debug(
			"Failed to fetch formats JSON",
			"url", url,
			"error", err,
		)
		return nil, err
	}

	e.logger.Debug(
		"Fetched formats JSON",
		"url", url,
		"elapsed", uformat.DurationFormat(elapsed),
	)

	err = e.formatCache.WriteByURL(url, dataJSON)
	if err != nil {
		e.logger.Error(
			"Failed to write formats to cache",
			"url", url,
			"error", err,
		)
		return nil, err
	}

	e.logger.Debug("Updated formats cache", "url", url)

	return dataJSON, nil
}

func (e *Executor) fetchInfoJSON(
	ctx context.Context,
	url string,
	opts ...idto.ExecutorOption,
) ([]byte, error) {
	// Prepare command arguments
	var args []string

	args = append(args, "--no-playlist", "--no-warnings", "--quiet")
	args = append(args, "--dump-json")

	// Add YouTube cookies if allowed in service options
	args = addCookiesToArgs(e.logger, args, opts...)

	args = append(args, url)

	// Execute the command to get video formats in JSON format
	out, err := utils.ExecCommandContext(ctx, e.ytDlpPath, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s command: %w", consts.YtDlpName, err)
	}

	outStr := string(out)

	// Find JSON start
	start := strings.Index(outStr, "{")
	if start == -1 {
		return nil, fmt.Errorf("no JSON found in %s output", consts.YtDlpName)
	}

	dataJSON := []byte(outStr[start:])

	var js json.RawMessage

	if err := json.Unmarshal(dataJSON, &js); err != nil {
		return nil, fmt.Errorf("invalid JSON output from %s: %w", consts.YtDlpName, err)
	}

	return dataJSON, nil
}
