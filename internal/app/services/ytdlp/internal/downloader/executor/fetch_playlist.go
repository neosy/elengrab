package executor

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/neosy/elengrab/internal/app/services/internal/utils"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

func (e *Executor) FetchPlaylist(
	ctx context.Context,
	playlistURL string,
	opts ...idto.ExecutorOption,
) (*idto.Playlist, error) {
	startTime := time.Now()
	extractItems, err := e.fetchAndExtractPlaylist(ctx, playlistURL, opts...)
	elapsed := time.Since(startTime)
	if err != nil {
		e.logger.Error(
			"Failed to fetch playlist",
			"url", playlistURL,
			"error", err,
		)
		return nil, fmt.Errorf("fetch and extract playlist error: %w", err)
	}

	e.logger.Debug(
		"Fetched playlist",
		"url", playlistURL,
		"elapsed", uformat.DurationFormat(elapsed),
	)

	items := make([]idto.PlaylistItem, 0, len(extractItems))

	for _, exItem := range extractItems {
		items = append(items, idto.PlaylistItem{
			URL:   exItem.URL,
			Title: exItem.Title,
		})
	}

	return &idto.Playlist{
		Items: items,
	}, nil
}

func (e *Executor) fetchAndExtractPlaylist(
	ctx context.Context,
	playlistURL string,
	opts ...idto.ExecutorOption,
) ([]idto.ExtractPlaylistItem, error) {
	// Prepare command arguments
	var args []string

	args = append(args, "--flat-playlist", "--no-warnings", "--dump-json")

	// Add cookies if allowed in service options
	args = addCookiesToArgs(e.logger, args, opts...)

	// args = append(args, "--load-info-json", e.formatCache.CacheFilePath(playlistURL))

	args = append(args, playlistURL)

	// Execute the command to get playlist in JSON format
	out, err := utils.ExecCommandContext(ctx, e.ytDlpPath, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s command: %w", consts.YtDlpName, err)
	}

	outStr := out.String()

	// Find JSON start
	found := strings.Contains(outStr, "{")
	if !found {
		return nil, fmt.Errorf("no JSON found in %s output", consts.YtDlpName)
	}

	items := make([]idto.ExtractPlaylistItem, 0)

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for scanner.Scan() {
		var item idto.ExtractPlaylistItem
		if err := json.Unmarshal(scanner.Bytes(), &item); err != nil {
			return nil, fmt.Errorf("invalid JSON output from %s: %w", consts.YtDlpName, err)
		}

		items = append(items, item)
	}

	return items, nil
}
