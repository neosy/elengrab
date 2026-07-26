package dto

type GetShareLinkResponse struct {
	DownloadID string `json:"itemId"`
	URL        string `json:"url"`
}
