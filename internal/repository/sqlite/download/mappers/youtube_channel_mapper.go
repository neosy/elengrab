package mappers

import (
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapYoutubeChannelDomainToEntity(channel *dmedia.YoutubeChannel) (*edownload.YoutubeChannel, error) {
	return &edownload.YoutubeChannel{
		ChannelID:   channel.ChannelID,
		ImageURL:    channel.ImageURL,
		ImageRaw:    channel.ImageRaw,
		ImageFormat: channel.ImageFormat,
	}, nil
}

func (m *Mappers) MapYoutubeChannelEntityToDomain(eChannel *edownload.YoutubeChannel) (*dmedia.YoutubeChannel, error) {
	return &dmedia.YoutubeChannel{
		ChannelID:   eChannel.ChannelID,
		ImageURL:    eChannel.ImageURL,
		ImageRaw:    eChannel.ImageRaw,
		ImageFormat: eChannel.ImageFormat,
		CreatedAt:   eChannel.CreatedAt,
		UpdatedAt:   eChannel.UpdatedAt,
	}, nil
}
