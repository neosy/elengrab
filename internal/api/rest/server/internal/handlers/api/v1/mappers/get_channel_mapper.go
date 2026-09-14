package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api/v1/dto"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
)

func (m *Mappers) MapChannelDomainToResponse(channel *dmedia.Channel) (*dto.GetChannelByIDResponse, error) {
	var imageURL, imageFormat string

	if channel.Image != nil {
		imageURL = channel.Image.URL
		imageFormat = channel.Image.Format.String()
	}

	return &dto.GetChannelByIDResponse{
		ChannelID:   channel.ChannelID.String(),
		ImageURL:    imageURL,
		ImageFormat: imageFormat,
	}, nil
}
