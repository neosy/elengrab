package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api/v1/dto"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

func (m *Mappers) MapChannelDomainToResponse(channel *dmedia.YoutubeChannel) (*dto.GetChannelByIDResponse, error) {
	return &dto.GetChannelByIDResponse{
		ChannelID:   channel.ChannelID,
		ImageURL:    channel.ImageURL,
		ImageFormat: channel.ImageFormat.String(),
	}, nil
}
