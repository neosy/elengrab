package searchindex

import (
	"context"
	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type SearchIndexPatchData struct {
	DownloadID uuid.UUID

	IsShort bool
	IsAudio bool
}

// CopyDataFrom copies all fields except DownloadID from src into the receiver.
func (d *SearchIndexPatchData) CopyDataFrom(src *SearchIndexPatchData) {
	if src == nil {
		return
	}

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
			if index == nil {
				return false
			}

			patchData := SearchIndexPatchData{
				DownloadID: index.DownloadID,

				IsShort: index.IsShort,
				IsAudio: index.IsAudio,
			}

			patch(&patchData)

			index.IsShort = patchData.IsShort
			index.IsAudio = patchData.IsAudio

			return true
		},
	)
}
