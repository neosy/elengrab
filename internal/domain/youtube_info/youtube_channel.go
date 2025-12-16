package dyoutube

import (
	"time"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type YoutubeChannel struct {
	// Unique ID for the channel
	ChannelID string

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

func (c *YoutubeChannel) InitFromResultChannelAvatar(chAvatar *ddownload.DownloadResultChannelAvatar) {
	if chAvatar != nil {
		c.ImageURL = chAvatar.ImageURL
		c.ImageRaw = chAvatar.ImageRAW
		c.ImageFormat = chAvatar.ImageFormat
	}
}
