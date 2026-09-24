package dmedia

import (
	"errors"
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/imgx"
)

type Channel struct {
	// Internal channel identifier
	ChannelID uuid.UUID

	// External channel identifier
	ExternalID string

	// Platform identifier
	Platform string

	// Channel URL
	ChannelURL string

	// Host from which the platform was detected
	Host string

	// Channel username
	Username string

	// Channel URL based on the username.
	UsernameURL string

	// Title of the channel
	Title string

	// Channel avatar image
	Image *dtypes.ChannelImage

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time
}

func NewChannelFromSource(source *dtypes.ChannelSource) *Channel {
	channel := &Channel{
		ChannelID: uuid.New(),
	}

	channel.InitFromSource(source)

	return channel
}

func (c *Channel) InitFromSource(channel *dtypes.ChannelSource) {
	if channel == nil {
		return
	}

	c.ExternalID = channel.ChannelID
	c.Platform = channel.Platform

	c.ChannelURL = channel.URL
	c.Host = channel.Host

	c.Username = channel.Username
	c.UsernameURL = channel.UsernameURL

	c.Title = channel.Title

	c.Image = channel.Image.Clone()
}

func (c *Channel) UpdateFromSource(channel *dtypes.ChannelSource) {
	if channel == nil {
		return
	}

	if channel.ChannelID != "" {
		c.ExternalID = channel.ChannelID
	}
	if channel.Platform != "" {
		c.Platform = channel.Platform
	}

	if channel.URL != "" {
		c.ChannelURL = channel.URL
	}
	if channel.Host != "" {
		c.Host = channel.Host
	}

	if channel.Username != "" {
		c.Username = channel.Username
	}
	if channel.UsernameURL != "" {
		c.UsernameURL = channel.UsernameURL
	}

	if channel.Title != "" {
		c.Title = channel.Title
	}

	if channel.Image != nil {
		c.Image = channel.Image.Clone()
	}
}

func (c *Channel) EqualSource(channel *dtypes.ChannelSource) bool {
	if channel == nil {
		return false
	}

	return c.ExternalID == channel.ChannelID &&
		c.Platform == channel.Platform &&
		c.ChannelURL == channel.URL &&
		c.Host == channel.Host &&
		c.Username == channel.Username &&
		c.UsernameURL == channel.UsernameURL &&
		c.Title == channel.Title &&
		c.Image.Equal(channel.Image)
}

func (c *Channel) Clone() *Channel {
	if c == nil {
		return nil
	}

	copyChannel := new(*c)

	copyChannel.Image = c.Image.Clone()

	return copyChannel
}

func (c *Channel) Validate() error {
	if c == nil {
		return errors.New("channel is nil")
	}

	if c.ChannelID == uuid.Nil {
		return errors.New("channel ID is empty")
	}

	if c.ExternalID == "" {
		return errors.New("channel external ID is empty")
	}

	if c.Platform == "" {
		return errors.New("channel platform is empty")
	}

	if c.ChannelURL == "" {
		return errors.New("channel URL is empty")
	}

	if c.Host == "" {
		return errors.New("channel host is empty")
	}

	return nil
}

func (c *Channel) IsValid() bool {
	return c.Validate() == nil
}

func (c *Channel) ImageData() *dtypes.ImageData {
	if c == nil || c.Image == nil || len(c.Image.Raw) == 0 {
		return nil
	}

	// decode
	size, _ := imgx.ImageSize(c.Image.Raw)

	return &dtypes.ImageData{
		URL:    c.Image.URL,
		Raw:    c.Image.Raw,
		Format: c.Image.Format,
		Width:  size.Width,
		Height: size.Height,
	}
}

func (c *Channel) ExternalChannelKey() dtypes.ExternalChannelKey {
	return dtypes.NewExternalChannelKey(c.ExternalID, c.Platform)
}

func (c *Channel) HasImage() bool {
	if c == nil || c.Image == nil {
		return false
	}

	return c.Image.IsValid()
}

func (c *Channel) PlatformType() dtypes.MediaPlatformType {
	platform, err := dtypes.ParseMediaPlatformType(c.Platform)
	if err != nil {
		return dtypes.MediaPlatformTypeNone
	}

	return platform
}
