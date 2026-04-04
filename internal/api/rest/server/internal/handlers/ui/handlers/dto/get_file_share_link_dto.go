package dto

type GetFileShareLinkResponse struct {
	FileID string `json:"fileId"`
	URL    string `json:"url"`
}
