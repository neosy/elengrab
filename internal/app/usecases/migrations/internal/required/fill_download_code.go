package required

import (
	"context"
	"errors"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

func (m *migrations) fillDownloadCode(ctx context.Context) error {
	const maxInsertAttempts = 3
	var downloadIDs []uuid.UUID

	err := m.Usecases().MediaDownload.IterateAllWithDeleted(ctx,
		func(d *ddownload.MediaDownload) error {
			if d == nil {
				return nil
			}

			downloadIDs = append(downloadIDs, d.DownloadID)

			return nil
		},
	)
	if err != nil {
		return err
	}

	if len(downloadIDs) == 0 {
		return nil
	}

	for _, downloadID := range downloadIDs {
		var err error
		for attempt := range maxInsertAttempts {
			err = m.Usecases().MediaDownload.Patch(ctx, nil, downloadID,
				func(d *ddownload.MediaDownload) error {
					if d == nil {
						return nil
					}

					d.DownloadCode = idcodec.EncodeUUIDShortBase64URL(d.DownloadID)

					return nil
				},
			)
			if err == nil || !errors.Is(err, apperrors.ErrIDConflict) {
				break
			}

			downloadID[15] = byte(attempt + 1)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
