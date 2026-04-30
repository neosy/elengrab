package imgx

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
)

type Size struct {
	Width  int
	Height int
}

// ImageSize decodes the image data and returns its dimensions (width and height).
func ImageSize(data []byte) (Size, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return Size{}, err
	}

	b := img.Bounds()

	return Size{
		Width:  b.Dx(),
		Height: b.Dy(),
	}, nil
}
