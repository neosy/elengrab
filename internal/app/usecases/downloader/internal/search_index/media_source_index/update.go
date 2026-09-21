package sourceindex

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaSourceIndex) IteratePatch(
	ctx context.Context,
	patch func(*ddownload.MediaSourceIndex) bool,
) error {
	return uc.indexRepo().Tx(ctx, func(ctx context.Context) error {
		var indexPatched []*ddownload.MediaSourceIndex

		err := uc.indexRepo().IterateAll(ctx, func(index *ddownload.MediaSourceIndex) error {
			if patch(index) {
				indexPatched = append(indexPatched, index)
			}
			return nil
		})
		if err != nil {
			return err
		}

		if len(indexPatched) == 0 {
			return nil
		}

		for _, index := range indexPatched {
			err := uc.indexRepo().Update(ctx, index)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
