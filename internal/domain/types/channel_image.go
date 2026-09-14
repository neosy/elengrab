package dtypes

import (
	"bytes"
	"errors"
)

type ChannelImage struct {
	// URL of the image channel avatar
	URL string

	// Raw image data (binary)
	Raw []byte

	// Format of the image (jpg, png, webp)
	Format ImageFormat
}

func (i *ChannelImage) Equal(other *ChannelImage) bool {
	if i == nil || other == nil {
		return i == other
	}

	return i.URL == other.URL &&
		bytes.Equal(i.Raw, other.Raw) &&
		i.Format == other.Format
}

func (image *ChannelImage) Clone() *ChannelImage {
	if image == nil {
		return nil
	}

	copyImage := new(*image)

	copyImage.Raw = bytes.Clone(image.Raw)

	return copyImage
}

func (i *ChannelImage) IsValid() bool {
	return i.Validate() == nil
}

func (i *ChannelImage) Validate() error {
	if i == nil {
		return errors.New("channel image is nil")
	}

	if i.URL == "" {
		return errors.New("channel image URL is empty")
	}

	if len(i.Raw) == 0 {
		return errors.New("channel image raw data is empty")
	}

	if !i.Format.Exists() {
		return errors.New("channel image format is invalid")
	}

	return nil
}
