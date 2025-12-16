package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api/v1/dto"
	dyoutube "github.com/neosy/elengrab/internal/domain/youtube_info"
)

func (m *Mappers) MapChannelDomainToResponse(channel *dyoutube.YoutubeChannel) (*dto.GetChannelByIDResponse, error) {
	return &dto.GetChannelByIDResponse{
		ChannelID:   channel.ChannelID,
		ImageURL:    channel.ImageURL,
		ImageFormat: channel.ImageFormat,
	}, nil
}
