package required

import (
	"context"
	"sync"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

var (
	fillSearchIndexChannelIDOnceSync sync.Once
)

func (m *migrations) fillSearchIndexChannelIDOnce(ctx context.Context) (bool, error) {
	var (
		ok  bool
		err error
	)

	fillSearchIndexChannelIDOnceSync.Do(func() {
		ok, err = m.fillSearchIndexChannelID(ctx)
	})

	return ok, err
}

func (m *migrations) fillSearchIndexChannelID(ctx context.Context) (bool, error) {
	channelIDsByDownloadID := make(map[uuid.UUID]uuid.UUID)

	err := m.Usecases().Downloader.MediaDownload().IterateAll(
		ctx, func(download *ddownload.MediaDownload) error {
			if download.ChannelID != uuid.Nil {
				channelIDsByDownloadID[download.DownloadID] = download.ChannelID
			}
			return nil
		})
	if err != nil {
		return false, err
	}

	err = m.Usecases().SearchIndex.IterateAllAndPatchChannelIDs(ctx, channelIDsByDownloadID)
	if err != nil {
		return false, err
	}

	return true, nil
}
