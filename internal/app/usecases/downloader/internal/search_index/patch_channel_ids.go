package searchindex

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *SearchIndex) IterateAllAndPatchChannelIDs(
	ctx context.Context,
	channelIDs map[uuid.UUID]uuid.UUID,
) error {
	return uc.sourceIndex.IteratePatch(ctx, func(index *ddownload.MediaSourceIndex) bool {
		newChannelID, exists := channelIDs[index.DownloadID]
		if !exists {
			return false
		}

		if newChannelID == index.ChannelID {
			return false
		}

		index.ChannelID = newChannelID

		return true
	})
}
