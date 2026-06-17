package mediadownload

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *MediaDownload) Patch(
	ctx context.Context,
	userID *uuid.UUID,
	downloadID uuid.UUID,
	mutate func(*ddownload.MediaDownload),
) error {
	download, err := uc.GetByDownloadID(ctx, userID, downloadID)
	if err != nil {
		return err
	}

	mutate(download)

	err = uc.Update(ctx, userID, download)
	if err != nil {
		return err
	}

	return nil
}

func (uc *MediaDownload) PatchMediaInfo(
	ctx context.Context,
	userID *uuid.UUID, downloadID uuid.UUID,
	mutate func(mediaInfo *dtypes.MediaInfo),
) error {
	download, err := uc.GetByDownloadID(ctx, userID, downloadID)
	if err != nil {
		return err
	}

	mutate(download.MediaInfo)

	err = uc.Update(ctx, nil, download)
	if err != nil {
		return err
	}

	return nil
}
