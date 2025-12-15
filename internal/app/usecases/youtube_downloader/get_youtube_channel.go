package ytdownloader

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
)

func (uc *YouTubeDownloader) GetYoutubeChannelInfo(ctx context.Context, channelID string) (*dto.GetYoutubeChannelInfoResponse, error) {
	channel, err := uc.ytChannel.FindByChannelId(ctx, channelID)
	if err != nil {
		uc.logger.Error("Failed get youtube channel info", "error", err)
		return nil, errorx.NewByErr(err, exceptionx.ERROR)
	}

	var resp *dto.GetYoutubeChannelInfoResponse
	if channel != nil {
		resp = uc.mappers.MapYoutubeChannelDomainToInfoResponse(channel)
	}

	return resp, nil
}
