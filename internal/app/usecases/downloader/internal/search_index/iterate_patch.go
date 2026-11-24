package searchindex

import (
	"context"
	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type SearchIndexPatchData struct {
	DownloadID uuid.UUID

	SourceType dtypes.MediaDownloadSourceType

	IsShort bool
	IsAudio bool
}

// NewSearchIndexPatchData creates a new SearchIndexPatchData from the given MediaSourceIndex.
func NewSearchIndexPatchData(src *ddownload.MediaSourceIndex) *SearchIndexPatchData {
	if src == nil {
		return &SearchIndexPatchData{}
	}

	return &SearchIndexPatchData{
		DownloadID: src.DownloadID,
		SourceType: src.SourceType,

		IsShort: src.IsShort,
		IsAudio: src.IsAudio,
	}
}

// InitFromMediaDownload initializes the receiver from the given download.
func (d *SearchIndexPatchData) InitFromMediaDownload(download *ddownload.MediaDownload) {
	d.DownloadID = download.DownloadID

	d.SourceType = download.SourceType

	d.IsShort = download.IsShort()
	d.IsAudio = download.IsAudio()
}

// CopyDataFrom copies all fields except DownloadID from src into the receiver.
func (d *SearchIndexPatchData) CopyDataFrom(src *SearchIndexPatchData) {
	if src == nil {
		return
	}

	d.SourceType = src.SourceType

	d.IsShort = src.IsShort
	d.IsAudio = src.IsAudio
}

func (uc *SearchIndex) IteratePatch(
	ctx context.Context,
	patch func(*SearchIndexPatchData) bool,
) error {
	return uc.sourceIndex.IteratePatch(
		ctx,
		func(index *ddownload.MediaSourceIndex) bool {
			if index == nil || patch == nil {
				return false
			}

			patchData := NewSearchIndexPatchData(index)

			if !patch(patchData) {
				return false
			}

			index.SourceType = patchData.SourceType

			index.IsShort = patchData.IsShort
			index.IsAudio = patchData.IsAudio

			return true
		},
	)
}
