package ytdlp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/yt_dlp/dto"
	dyoutubeinfo "github.com/neosy/elengrab/internal/domain/youtube_info"
)

// GetFormats retrieves and parses video formats for the given URL.
func (y *YTDlp) GetFormats(ctx context.Context, url string) (*dyoutubeinfo.YouTubeInfo, error) {
	info, err := y.getFormats(ctx, url)
	if err != nil {
		return nil, err
	}

	return y.mappers.VideoInfoToDomain(info), nil
}

func (y *YTDlp) fetchFormatsJSON(
	ctx context.Context,
	url string,
) ([]byte, error) {
	out, err := y.execCommandContext(
		ctx, y.ytDlpPath,
		"--no-playlist", "--no-warnings", "--quiet",
		"--dump-json",
		url,
	)
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

func (y *YTDlp) fetchAndCacheFormatsJSON(
	ctx context.Context,
	url string,
) ([]byte, error) {
	dataJSON, err := y.fetchFormatsJSON(ctx, url)
	if err != nil {
		y.formatCache.deleteByURL(url)
		return nil, err
	}

	err = y.formatCache.writeByURL(url, dataJSON)
	if err != nil {
		return nil, err
	}

	return dataJSON, nil
}

func (y *YTDlp) updateCache(
	ctx context.Context,
	url string,
) error {
	valid, err := y.formatCache.isTTLValidByURL(url)
	if err != nil {
		return err
	}

	if valid {
		return nil
	}

	return y.updateCacheForce(ctx, url)
}

func (y *YTDlp) updateCacheForce(
	ctx context.Context,
	url string,
) error {
	_, err := y.fetchAndCacheFormatsJSON(ctx, url)
	if err != nil {
		return err
	}

	return nil
}

func (y *YTDlp) getFormatsJSON(
	ctx context.Context,
	url string,
) ([]byte, error) {
	dataJSON, err := y.formatCache.loadByURL(url)
	if err != nil {
		return nil, err
	}

	if dataJSON != nil {
		return dataJSON, nil
	}

	dataJSON, err = y.fetchAndCacheFormatsJSON(ctx, url)
	if err != nil {
		return nil, err
	}

	return dataJSON, nil
}

func (y *YTDlp) getFormats(
	ctx context.Context,
	url string,
) (*idto.YouTubeInfo, error) {
	dataJSON, err := y.getFormatsJSON(ctx, url)
	if err != nil {
		return nil, err
	}

	var info = &idto.YouTubeInfo{}
	err = json.NewDecoder(bytes.NewReader(dataJSON)).Decode(info)
	if err != nil {
		y.formatCache.deleteByURL(url)
		return nil, fmt.Errorf("json decode error: %v", err)
	}

	return info, nil
}
