package ddownload

type DownloadChannel struct {
	URL    string
	Title  string
	Avatar *DownloadChannelAvatar
}

type DownloadChannelAvatar struct {
	ImageURL    string
	ImageRAW    []byte
	ImageFormat string
}

func (c *DownloadChannel) Copy() *DownloadChannel {
	if c == nil {
		return nil
	}

	channel := *c
	channel.Avatar = c.Avatar.Copy()

	return &channel
}

func (a *DownloadChannelAvatar) Copy() *DownloadChannelAvatar {
	if a == nil {
		return nil
	}

	avatar := *a
	// Deep copy of ImageRAW
	avatar.ImageRAW = make([]byte, len(a.ImageRAW))
	copy(avatar.ImageRAW, a.ImageRAW)

	return &avatar
}
