package dtypes

type ImageData struct {
	// URL of the image
	URL string

	// Format of the image (jpg, png, webp)
	Format ImageFormat

	// Dimensions of the image
	Width  int
	Height int

	// Raw image data (binary)
	Raw []byte
}

func (d *ImageData) Copy() *ImageData {
	if d == nil {
		return nil
	}

	data := *d

	raw := make([]byte, len(d.Raw))
	copy(raw, d.Raw)

	return &data
}

func (d *ImageData) FullURL(baseURL string) string {
	return baseURL + d.URL
}
