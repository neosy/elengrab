package dto

type GetFileShareLinkResponse struct {
	FileID string `json:"downloadID"`
	URL    string `json:"url"`
}
