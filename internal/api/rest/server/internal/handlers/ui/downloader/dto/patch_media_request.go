package dto

type PatchMediaByDownloadIDRequest struct {
	Title       string `json:"title" validate:"required,max=100"`
	Description string `json:"description" validate:"max=5000"`
	Visibility  string `json:"visibility" validate:"required,visibility"`
}
