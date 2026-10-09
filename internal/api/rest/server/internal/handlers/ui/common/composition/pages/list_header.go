package pages

type ListHeader struct {
	Title       string `json:"title"`
	Description string `json:"description"`

	ImageURL string `json:"imageUrl"`

	Show bool `json:"show"`
}
