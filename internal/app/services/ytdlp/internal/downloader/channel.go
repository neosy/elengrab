package downloader

import (
	"context"
	"net/url"
	"time"

	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/helper"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	"github.com/neosy/elengrab/internal/app/utils/siteimage/channels"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
)

func (d *Downloader) buildChannel(mediaURL string, extractInfo *idto.ExtractInfo) *dtypes.ChannelSource {
	if extractInfo == nil {
		return nil
	}

	if extractInfo.ChannelURL == "" && mediaURL == "" {
		return nil
	}

	var platform string
	if extractInfo.ChannelURL != "" {
		platform = hostdetect.DetectPlatformName(extractInfo.ChannelURL)
	} else {
		platform = hostdetect.DetectPlatformName(mediaURL)
	}

	if platform == "" {
		return nil
	}

	var channelHost string
	if extractInfo.ChannelURL != "" {
		parsedURL, err := url.Parse(extractInfo.ChannelURL)
		if err == nil {
			channelHost = parsedURL.Host
		}
	} else {
		parsedURL, err := url.Parse(mediaURL)
		if err == nil {
			channelHost = parsedURL.Host
		}
	}

	channelSource := &dtypes.ChannelSource{
		ChannelID: extractInfo.ChannelID,
		Platform:  platform,

		URL:  extractInfo.ChannelURL,
		Host: channelHost,

		Username:    extractInfo.ParsedChannel.Username,
		UsernameURL: extractInfo.ParsedChannel.UsernameURL,

		Title: extractInfo.ParsedChannel.Title,
	}

	if channelSource.IsValid() {
		return channelSource
	}

	return nil
}

func (d *Downloader) fetchChannelImage(
	ctx context.Context,
	channelURL string,
	opts ...channels.FetchOption,
) (*dtypes.ChannelImage, error) {
	if channelURL == "" {
		return nil, nil
	}

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

func (d *Downloader) buildChannelFetchOptions(channelURL string, cookieFilePath string) []channels.FetchOption {
	var opts []channels.FetchOption
	opts = append(opts, idto.DefaultRequestOptions().ChannelFetchOption())

	if cookieFilePath != "" && hostdetect.Instagram(channelURL) {
		cookies, _ := helper.ParseCookiesFile(cookieFilePath)
		if len(cookies) > 0 {
			opts = append(opts, channels.FetchOptionsWithCookies(cookies))
		}
	}

	return opts
}
