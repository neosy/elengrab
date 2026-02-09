package mappers

import (
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	emedia "github.com/neosy/elengrab/internal/repository/sqlite/media/entity"
)

func (m *Mappers) MapYoutubeChannelDomainToEntity(channel *dmedia.YoutubeChannel) (*emedia.YoutubeChannel, error) {
	return &emedia.YoutubeChannel{
		ChannelID:   channel.ChannelID,
		ImageURL:    channel.ImageURL,
		ImageRaw:    channel.ImageRaw,
		ImageFormat: channel.ImageFormat,
	}, nil
}

func (m *Mappers) MapYoutubeChannelEntityToDomain(eChannel *emedia.YoutubeChannel) (*dmedia.YoutubeChannel, error) {
	return &dmedia.YoutubeChannel{
		ChannelID:   eChannel.ChannelID,
		ImageURL:    eChannel.ImageURL,
		ImageRaw:    eChannel.ImageRaw,
		ImageFormat: eChannel.ImageFormat,
		CreatedAt:   eChannel.CreatedAt,
		UpdatedAt:   eChannel.UpdatedAt,
	}, nil
}
