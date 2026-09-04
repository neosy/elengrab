package sourceindex

import (
	"context"

	"github.com/google/uuid"
)

func (uc *MediaSourceIndex) UpdateViews(ctx context.Context, downloadID uuid.UUID, views uint32) error {
	return uc.Tx(ctx, func(ctx context.Context) error {
		index, err := uc.FindByDownload(ctx, downloadID)
		if err != nil {
			return err
		}

		if index == nil {
			return nil
		}

		if index.Views == views {
			return nil
		}

		index.Views = views

		err = uc.Update(ctx, index)
		if err != nil {
			return err
		}

		return nil
	})
}
