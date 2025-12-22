package idto

type AvatarSource struct {
	URL    string `json:"url"`
	Format string
	Raw    []byte
}
