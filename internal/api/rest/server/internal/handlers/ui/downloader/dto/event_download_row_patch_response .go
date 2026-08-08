package dto

import "html/template"

type EventDownloadRowPatchResponse struct {
	DownloadID string `json:"itemId"`

	Title        string              `json:"title,omitempty"`
	Visibility   *VisibilityResponse `json:"visibility,omitempty"`
	HasShareLink *bool               `json:"hasShareLink,omitempty"`

	WatchPercent string             `json:"watchPercent,omitempty"`
	Watched      *bool              `json:"watched,omitempty"`
	ViewCount    *ViewCountResponse `json:"viewCount,omitempty"`
}

type VisibilityResponse struct {
	Visible bool          `json:"visible"`
	Value   string        `json:"value"`
	Label   string        `json:"label"`
	Icon    template.HTML `json:"icon"`
}

type ViewCountResponse struct {
	Count uint32 `json:"count"`
	Text  string `json:"text"`
}

func (r *EventDownloadRowPatchResponse) HasChanges() bool {
	return r.Title != "" ||
		r.Visibility != nil ||
		r.HasShareLink != nil ||
		r.WatchPercent != "" ||
		r.Watched != nil ||
		r.ViewCount != nil
}
