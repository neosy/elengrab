package required

import (
	"context"
	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (m *migrations) updateSearchIndex(ctx context.Context) error {
	patchDataByDownloadID := make(map[uuid.UUID]*downloader.SearchIndexPatchData)

	err := m.Usecases().Downloader.MediaDownload().IterateAll(
		ctx, func(download *ddownload.MediaDownload) error {
			patchData := downloader.NewSearchIndexPatchData(nil)
			patchData.InitFromMediaDownload(download)

			patchDataByDownloadID[download.DownloadID] = patchData

			return nil
		})
	if err != nil {
		return err
	}

	err = m.Usecases().SearchIndex.IteratePatch(
		ctx,
		func(patchData *downloader.SearchIndexPatchData) bool {
			if patchData == nil {
				return false
			}

			data, exists := patchDataByDownloadID[patchData.DownloadID]
			if !exists {
				return false
			}

			patchData.CopyDataFrom(data)

			return true
		},
	)
	if err != nil {
		return err
	}

	return nil
}
