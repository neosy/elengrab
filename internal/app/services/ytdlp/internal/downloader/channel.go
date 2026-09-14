package downloader

import (
	"context"
	"net/url"
	"time"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	"github.com/neosy/elengrab/internal/app/utils/siteimage/channels"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

func (d *Downloader) fetchAndBuildChannel(
	ctx context.Context,
	meta *idto.DownloadMeta,
	opts ...channels.FetchOption,
) *dtypes.ChannelSource {
	if meta == nil {
		return nil
	}

	if meta.ChannelURL == "" {
		return nil
	}

	parsedURL, err := url.Parse(meta.ChannelURL)
	if err != nil {
		return nil
	}

	host := parsedURL.Host

	var (
		elapsed time.Duration
	)

	channel := &dtypes.ChannelSource{
		URL: meta.ChannelURL,

		Platform: hostdetect.DetectPlatformType(meta.ChannelURL).String(),
		Host:     host,

		ChannelID: meta.ChannelID,
		Title:     meta.ChannelTitle,
	}

	startTime := time.Now()
	channelImages, err := channels.FetchImages(ctx, meta.ChannelURL, opts...)
	elapsed = time.Since(startTime)
	if err != nil {
		d.logger.Debug("Failed to get channel avatar", "channelURL", meta.ChannelURL, "error", err)
	}

	if len(channelImages) == 0 {
		d.logger.Debug("Channel image not found", "channelURL", meta.ChannelURL)
		return channel
	}

	d.logger.Info(
		"Channel avatar fetched",
		"host", hostdetect.DetectPlatformType(meta.ChannelURL).String(),
		"channelURL", meta.ChannelURL,
		"elapsed", uformat.DurationFormat(elapsed),
	)

	channel.Image = new(dtypes.ChannelImage(channelImages[0]))

	return channel
}
