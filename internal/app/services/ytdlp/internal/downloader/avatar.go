package downloader

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

func (d *Downloader) fetchChannelAvatarAsync(
	wg *sync.WaitGroup,
	meta *idto.DownloadMeta,
	onChannelDone func(*ddownload.DownloadChannel),
) {
	if meta.ChannelURL == "" {
		return
	}

	wg.Go(func() {
		startTime := time.Now()
		avatarSources, err := d.fetchChannelAvatar(meta.ChannelURL)
		elapsed := time.Since(startTime)
		if err != nil {
			d.logger.Debug("Failed to get channel avatar", "channelURL", meta.ChannelURL, "error", err)
			return
		}
		if len(avatarSources) == 0 {
			d.logger.Debug("Avatar not found", "channelURL", meta.ChannelURL)
			return
		}

		var channel *ddownload.DownloadChannel
		src := avatarSources[0]
		if len(src.Raw) == 0 {
			d.logger.Debug("Avatar image not found", "channelURL", meta.ChannelURL)
			return
		}

		d.logger.Info(
			"YouTube channel avatar fetched",
			"channelURL", meta.ChannelURL,
			"elapsed", uformat.DurationFormat(elapsed),
		)

		channel = &ddownload.DownloadChannel{
			URL:   meta.ChannelURL,
			Title: meta.ChannelTitle,
			Avatar: &ddownload.DownloadChannelAvatar{
				ImageURL:    src.URL,
				ImageRAW:    src.Raw,
				ImageFormat: src.Format,
			},
		}
		onChannelDone(channel)
	})
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
		return nil, fmt.Errorf("failed to get html: %w", err)
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
