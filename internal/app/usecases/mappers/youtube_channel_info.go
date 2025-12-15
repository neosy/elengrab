package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (m *Mappers) MapYoutubeChannelDomainToInfoResponse(channel *ddownload.YoutubeChannel) *dto.GetYoutubeChannelInfoResponse {
	return &dto.GetYoutubeChannelInfoResponse{
		ChannelID:   channel.ChannelID,
		ImageURL:    channel.ImageURL,
		ImageRaw:    channel.ImageRaw,
		ImageFormat: channel.ImageFormat,
	}
}
