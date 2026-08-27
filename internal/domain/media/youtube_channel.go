package dmedia

import (
	"bytes"
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/imgx"
)

type YoutubeChannel struct {
	// Unique ID for the channel
	ChannelID string

	// Site URL
	ChannelURL string

	// Title of the channel
	ChannelTitle string

	// URL of the channel avatar
	ImageURL string

	// Raw image data (binary)
	ImageRaw []byte

	// Format of the image (jpg, png, webp)
	ImageFormat dtypes.ImageFormat

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time
}

func (c *YoutubeChannel) InitFromChannel(channel *dtypes.Channel) {
	if channel != nil {
		c.ChannelURL = channel.URL
		c.ChannelTitle = channel.Title
		if channel.Avatar != nil {
			c.ImageURL = channel.Avatar.ImageURL
			c.ImageRaw = channel.Avatar.ImageRAW
			c.ImageFormat = channel.Avatar.ImageFormat
		}
	}
}

func (c *YoutubeChannel) Copy() *YoutubeChannel {
	if c == nil {
		return nil
	}

	copy := *c

	copy.ImageRaw = bytes.Clone(c.ImageRaw)

	return &copy
}

func (c *YoutubeChannel) ImageData() *dtypes.ImageData {
	if c == nil {
		return nil
	}

	if len(c.ImageRaw) == 0 {
		return nil
	}

	// decode
	size, _ := imgx.ImageSize(c.ImageRaw)

	return &dtypes.ImageData{
		URL:    c.ImageURL,
		Format: c.ImageFormat,
		Raw:    c.ImageRaw,
		Width:  size.Width,
		Height: size.Height,
	}
}
