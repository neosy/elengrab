package dto

type GetDownloadShareLinkResponse struct {
	DownloadID string `json:"itemId"`
	URL        string `json:"url"`
}
