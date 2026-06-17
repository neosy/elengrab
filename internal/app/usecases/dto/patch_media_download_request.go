package dto

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

type PatchMediaDownloadRequest struct {
	DownloadID uuid.UUID

	MediaTitle       *string
	MediaDescription *string
	Visibility       *dtypes.MediaVisibility
}

func (r *PatchMediaDownloadRequest) Normalize() {
	if r.MediaTitle != nil {
		var title string
		title = strings.TrimSpace(*r.MediaTitle)
		if title == "" {
			r.MediaTitle = nil
		} else {
			r.MediaTitle = new(title)
		}
	}

	if r.MediaDescription != nil {
		var description string
		description = strings.TrimSpace(*r.MediaDescription)
		if description == "" {
			r.MediaDescription = nil
		} else {
			r.MediaDescription = new(description)
		}
	}
}

func (r *PatchMediaDownloadRequest) Validate() error {
	if r.MediaTitle != nil && *r.MediaTitle == "" {
		return errorx.NewHTTPMessage("Title is required", http.StatusBadRequest)
	}

	return nil
}
