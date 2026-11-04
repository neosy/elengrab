package dto

import (
	dconsts "github.com/neosy/elengrab/internal/domain/consts"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

type QueueMediaDownloadListRequest struct {
	Items   []QueueMediaDownloadListRequestItem `json:"items" validate:"min=1"`
	Options MediaDownloadOptions                `json:"options"`
}

func (r *QueueMediaDownloadListRequest) Validate() error {
	if r == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	if len(r.Items) == 0 {
		return errorx.NewMessage("Media download list cannot be empty")
	}

	if len(r.Items) > dconsts.MediaDownloadListMaxLen {
		return errorx.NewMessage("Media download list exceeds maximum length")
	}

	return nil
}

type QueueMediaDownloadListResponse struct {
	Items []QueueMediaDownloadListResponseItem `json:"items"`
	Count int                                  `json:"count"`
}

type QueueMediaDownloadListRequestItem struct {
	URL   string `json:"url" validate:"required,url,max_bytes=2048"`
	Title string `json:"title" validate:"max=100"`
}

type QueueMediaDownloadListResponseItem struct {
	DownloadID string `json:"itemId"`

	URL   string `json:"url"`
	Title string `json:"title"`
}

type MediaDownloadOptions struct {
	QualityCodec      string `json:"qualityCodec" validate:"required"`
	QualityResolution string `json:"qualityResolution" validate:"required"`
	Format            string `json:"format" validate:"required"`
}
