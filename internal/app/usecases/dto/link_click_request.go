package dto

type ShortLinkClickRequest struct {
	ShortURL string

	ClickedBy *string

	IPAddress string
	UserAgent *string
	Referrer  *string
}
