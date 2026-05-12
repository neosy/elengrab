package mediadownload

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *MediaDownload) Patch(ctx context.Context, userID *uuid.UUID, downloadID uuid.UUID, patch *dto.MediaDownloadInfoPatch) error {
	file, err := uc.GetByDownloadID(ctx, userID, downloadID)
	if err != nil {
		return err
	}

	dto.PatchToMediaDownloadDomain(patch, file)

	err = uc.Update(ctx, file)
	if err != nil {
		return err
	}

	return nil
}

func (uc *MediaDownload) PatchRecord(
	ctx context.Context,
	userID *uuid.UUID, downloadID uuid.UUID,
	patch func(download *ddownload.MediaDownload)) error {
	download, err := uc.GetByDownloadID(ctx, userID, downloadID)
	if err != nil {
		return err
	}

	patch(download)

	err = uc.Update(ctx, download)
	if err != nil {
		return err
	}

	return nil
}

func (uc *MediaDownload) PatchMediaInfo(
	ctx context.Context,
	userID *uuid.UUID, downloadID uuid.UUID,
	patchMediaInfo func(mediaInfo *dtypes.MediaInfo)) error {
	download, err := uc.GetByDownloadID(ctx, userID, downloadID)
	if err != nil {
		return err
	}

	patchMediaInfo(download.MediaInfo)

	err = uc.Update(ctx, download)
	if err != nil {
		return err
	}

	return nil
}
