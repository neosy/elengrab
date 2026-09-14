package channels

type ImageSource struct {
	URL    string `json:"url"`
	Format string
	Raw    []byte
}
