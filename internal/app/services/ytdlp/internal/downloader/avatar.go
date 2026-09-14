package downloader

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/consts"
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

func (d *Downloader) fetchAndBuildChannel(meta *idto.DownloadMeta) *dtypes.ChannelSource {
	if meta.ChannelURL == "" {
		return nil
	}

	parsedURL, err := url.Parse(meta.ChannelURL)
	if err != nil {
		return nil
	}

	host := parsedURL.Host

	var (
		avatarSources []idto.AvatarSource
		elapsed       time.Duration
	)

	channel := &dtypes.ChannelSource{
		URL: meta.ChannelURL,

		Platform: hostdetect.DetectPlatformType(meta.ChannelURL).String(),
		Host:     host,

		ChannelID: meta.ChannelID,
		Title:     meta.ChannelTitle,
	}

	if hostdetect.YouTube(meta.ChannelURL) {
		var err error
		startTime := time.Now()
		avatarSources, err = d.fetchYoutubeChannelAvatar(meta.ChannelURL)
		elapsed = time.Since(startTime)
		if err != nil {
			d.logger.Debug("Failed to get channel avatar", "channelURL", meta.ChannelURL, "error", err)
		}
	}

	if len(avatarSources) == 0 {
		d.logger.Debug("Avatar image not found", "channelURL", meta.ChannelURL)
		return channel
	}

	avatarSource := &avatarSources[0]

	d.logger.Info(
		"Channel avatar fetched",
		"host", hostdetect.DetectPlatformType(meta.ChannelURL).String(),
		"channelURL", meta.ChannelURL,
		"elapsed", uformat.DurationFormat(elapsed),
	)

	imageFormat, err := dtypes.ParseImageFormat(avatarSource.Format)
	if err != nil {
		d.logger.Warn(
			"Failed to parse image format",
			"format", avatarSource.Format,
			"error", err,
		)
		return channel
	}

	channel.Image = &dtypes.ChannelImage{
		URL:    avatarSource.URL,
		Raw:    avatarSource.Raw,
		Format: imageFormat,
	}

	return channel
}

// fetchChannelAvatar fetches the HTML of a YouTube channel page,
// extracts the avatar JSON block, and returns all avatar URLs.
// url: full URL of the YouTube channel
func (d *Downloader) fetchYoutubeChannelAvatar(url string) ([]idto.AvatarSource, error) {
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

	var sources, loadedSources []idto.AvatarSource
	if err := json.Unmarshal([]byte(jsonArray), &sources); err != nil {
		return nil, err
	}

	for _, src := range sources {
		raw, format, err := nfasthttp.GetImage(
			src.URL,
			nfasthttp.ClientOptionWithreadBufferSize(64*1024),
			nfasthttp.ClientOptionWithTimeout(consts.ChannelAvatarTimeout),
		)
		if err != nil || len(raw) == 0 {
			continue
		}

		loadedSources = append(loadedSources,
			idto.AvatarSource{
				URL:    src.URL,
				Raw:    raw,
				Format: format,
			},
		)
	}

	return loadedSources, nil
}
