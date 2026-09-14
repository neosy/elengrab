package types

type ImageSource struct {
	URL    string `json:"url"`
	Format string
	Raw    []byte
}
