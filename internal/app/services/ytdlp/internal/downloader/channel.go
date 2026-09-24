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

func (d *Downloader) buildChannel(extractInfo *idto.ExtractInfo) *dtypes.ChannelSource {
	if extractInfo == nil || extractInfo.ChannelURL == "" {
		return nil
	}

	parsedURL, err := url.Parse(extractInfo.ChannelURL)
	if err != nil {
		return nil
	}

	return &dtypes.ChannelSource{
		ChannelID: extractInfo.ChannelID,
		Platform:  hostdetect.DetectPlatformName(extractInfo.ChannelURL),

		URL:  extractInfo.ChannelURL,
		Host: parsedURL.Host,

		Username:    extractInfo.ParsedChannel.Username,
		UsernameURL: extractInfo.ParsedChannel.UsernameURL,

		Title: extractInfo.ParsedChannel.Title,
	}
}

func (d *Downloader) fetchChannelImage(
	ctx context.Context,
	channelURL string,
	opts ...channels.FetchOption,
) (*dtypes.ChannelImage, error) {
	var elapsed time.Duration

	startTime := time.Now()
	channelImages, err := channels.FetchImages(ctx, channelURL, opts...)
	elapsed = time.Since(startTime)
	if err != nil {
		d.logger.Debug(
			"Failed to get channel avatar",
			"channelURL", channelURL,
			"error", err,
		)
		return nil, err
	}

	if len(channelImages) == 0 {
		d.logger.Debug("Channel image not found", "channelURL", channelURL)
		return nil, nil
	}

	d.logger.Info(
		"Channel avatar fetched",
		"channelURL", channelURL,
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return new(dtypes.ChannelImage(channelImages[0])), nil
}
