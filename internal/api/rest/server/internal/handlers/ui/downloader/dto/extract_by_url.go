package dto

type ExtractByURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type ExtractByURLResponse struct {
	Items []ExtractByURLItem `json:"items"`
}

type ExtractByURLItem struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}
