package dto

type ExtractByURLResponse struct {
	Items []ExtractByURLItem
}

type ExtractByURLItem struct {
	URL   string
	Title string
}
