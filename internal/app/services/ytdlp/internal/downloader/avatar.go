package downloader

import (
	"encoding/json"
	"time"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

func (d *Downloader) fetchAndBuildChannelAvatar(meta *idto.DownloadMeta) *dtypes.Channel {
	if meta.ChannelURL == "" {
		return nil
	}

	startTime := time.Now()
	avatarSources, err := d.fetchChannelAvatar(meta.ChannelURL)
	elapsed := time.Since(startTime)
	if err != nil {
		d.logger.Debug("Failed to get channel avatar", "channelURL", meta.ChannelURL, "error", err)
		return nil
	}
	if len(avatarSources) == 0 {
		d.logger.Debug("Avatar not found", "channelURL", meta.ChannelURL)
		return nil
	}

	src := avatarSources[0]
	if len(src.Raw) == 0 {
		d.logger.Debug("Avatar image not found", "channelURL", meta.ChannelURL)
		return nil
	}

	d.logger.Info(
		"YouTube channel avatar fetched",
		"channelURL", meta.ChannelURL,
		"elapsed", uformat.DurationFormat(elapsed),
	)

	imageFormat, err := dtypes.ParseImageFormat(src.Format)
	if err != nil {
		d.logger.Warn(
			"Failed to parse image format",
			"format", src.Format,
			"error", err,
		)
		return nil
	}

	return &dtypes.Channel{
		URL:   meta.ChannelURL,
		Title: meta.ChannelTitle,
		Avatar: &dtypes.ChannelAvatar{
			ImageURL:    src.URL,
			ImageRAW:    src.Raw,
			ImageFormat: imageFormat,
		},
	}
}

// fetchChannelAvatar fetches the HTML of a YouTube channel page,
// extracts the avatar JSON block, and returns all avatar URLs.
// url: full URL of the YouTube channel
func (d *Downloader) fetchChannelAvatar(url string) ([]idto.AvatarSource, error) {
	body, err := nfasthttp.GetHTML(
		url,
		nfasthttp.ClientOptionWithreadBufferSize(64*1024),
		nfasthttp.ClientOptionWithTimeout(consts.ChannelAvatarTimeout),
	)
	if err != nil {
		return nil, errorx.Errorf("failed to get html: %w", err)
	}

	html := string(body)

	// Full key path to the avatar sources inside decoratedAvatarViewModel
	fullKey := `"decoratedAvatarViewModel":{"avatar":{"avatarViewModel":{"image":{"sources":[`
	jsonArray, err := helper.ExtractJSONArray(html, fullKey)
	if err != nil {
		return nil, err
	}

	var sources []idto.AvatarSource
	if err := json.Unmarshal([]byte(jsonArray), &sources); err != nil {
		return nil, err
	}

	for i, src := range sources {
		raw, format, err := nfasthttp.GetImage(
			src.URL,
			nfasthttp.ClientOptionWithreadBufferSize(64*1024),
			nfasthttp.ClientOptionWithTimeout(consts.ChannelAvatarTimeout),
		)
		if err != nil {
			continue
		}
		sources[i].Raw = raw
		sources[i].Format = format
	}

	return sources, nil
}
