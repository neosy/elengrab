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

func (m *migrations) fillSearchIndexChannelIDOnce(ctx context.Context) error {
	var err error

	fillSearchIndexChannelIDOnceSync.Do(func() {
		err = m.fillSearchIndexChannelID(ctx)
	})

	return err
}

func (m *migrations) fillSearchIndexChannelID(ctx context.Context) error {
	channelIDsByDownloadID := make(map[uuid.UUID]uuid.UUID)

	err := m.Usecases().Downloader.MediaDownload().IterateAll(
		ctx, func(download *ddownload.MediaDownload) error {
			if download.ChannelID != uuid.Nil {
				channelIDsByDownloadID[download.DownloadID] = download.ChannelID
			}
			return nil
		})
	if err != nil {
		return err
	}

	err = m.Usecases().SearchIndex.IterateAllAndPatchChannelIDs(ctx, channelIDsByDownloadID)
	if err != nil {
		return err
	}

	return nil
}
