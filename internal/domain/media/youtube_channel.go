package dmedia

import (
	"time"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
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
	ImageFormat string

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time
}

func (c *YoutubeChannel) InitFromChannel(channel *ddownload.DownloadChannel) {
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

func (src *YoutubeChannel) Copy() *YoutubeChannel {
	if src == nil {
		return nil
	}

	copy := uptr.Copy(src)

	if len(src.ImageRaw) > 1 {
		copy.ImageRaw = append([]byte{}, src.ImageRaw...)
	}

	return copy
}
