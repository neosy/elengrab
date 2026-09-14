package mappers

import (
	"bytes"

	"github.com/google/uuid"

	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	emedia "github.com/neosy/elengrab/internal/repository/sqlite/media/entity"
)

func (m *Mappers) MapChannelDomainToEntity(channel *dmedia.Channel) (*emedia.Channel, error) {
	var (
		imageURL    *string
		imageRAW    []byte
		imageFormat *string
	)
	if channel.Image != nil {
		if err := channel.Image.Validate(); err != nil {
			return nil, err
		}

		imageURL = &channel.Image.URL
		imageRAW = bytes.Clone(channel.Image.Raw)
		imageFormat = new(channel.Image.Format.String())
	}

	return &emedia.Channel{
		ChannelID:   channel.ChannelID.String(),
		Platform:    channel.Platform,
		Host:        channel.Host,
		ExternalID:  channel.ExternalID,
		ChannelURL:  channel.ChannelURL,
		Title:       channel.Title,
		ImageURL:    imageURL,
		ImageRaw:    imageRAW,
		ImageFormat: imageFormat,
	}, nil
}

func (m *Mappers) MapChannelEntityToDomain(eChannel *emedia.Channel) (*dmedia.Channel, error) {
	channelID, err := uuid.Parse(eChannel.ChannelID)
	if err != nil {
		return nil, err
	}

	image := &dtypes.ChannelImage{
		Raw: eChannel.ImageRaw,
	}

	if eChannel.ImageURL != nil {
		image.URL = *eChannel.ImageURL
	}

	if eChannel.ImageFormat != nil {
		imageFormat, err := dtypes.ParseImageFormat(*eChannel.ImageFormat)
		if err != nil {
			return nil, err
		}
		image.Format = imageFormat
	}

	if !image.IsValid() {
		image = nil
	}

	return &dmedia.Channel{
		ChannelID:  channelID,
		Platform:   eChannel.Platform,
		Host:       eChannel.Host,
		ExternalID: eChannel.ExternalID,
		ChannelURL: eChannel.ChannelURL,
		Title:      eChannel.Title,
		Image:      image,
		CreatedAt:  eChannel.CreatedAt,
		UpdatedAt:  eChannel.UpdatedAt,
	}, nil
}
