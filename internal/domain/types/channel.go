package dtypes

type Channel struct {
	URL    string
	Title  string
	Avatar *ChannelAvatar
}

type ChannelAvatar struct {
	ImageURL    string
	ImageRAW    []byte
	ImageFormat ImageFormat
}

func (c *Channel) Copy() *Channel {
	if c == nil {
		return nil
	}

	channel := *c
	channel.Avatar = c.Avatar.Copy()

	return &channel
}

func (a *ChannelAvatar) Copy() *ChannelAvatar {
	if a == nil {
		return nil
	}

	avatar := *a
	// Deep copy of ImageRAW
	avatar.ImageRAW = make([]byte, len(a.ImageRAW))
	copy(avatar.ImageRAW, a.ImageRAW)

	return &avatar
}
