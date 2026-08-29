package downloader

import (
	"context"

	"github.com/google/uuid"
	mediadownload "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type MediaDownload interface {
	Patch(
		ctx context.Context,
		userID *uuid.UUID,
		downloadID uuid.UUID,
		mutate func(*ddownload.MediaDownload) error,
	) error
	PatchMediaInfo(
		ctx context.Context,
		userID *uuid.UUID, downloadID uuid.UUID,
		mutate func(mediaInfo *dtypes.MediaInfo),
	) error

	FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaDownload, error)
	GetAllFullNames(ctx context.Context) (map[string]struct{}, error)
	GetAllFullNamesWithDeleted(ctx context.Context) (map[string]struct{}, error)
	IterateGetAll(ctx context.Context, fn func(*ddownload.MediaDownload) error) error
	IterateGetAllWithDeleted(ctx context.Context, fn func(*ddownload.MediaDownload) error) error
}

type mediaDownload struct {
	*mediadownload.MediaDownload
}

func (d *mediaDownload) FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaDownload, error) {
	return d.MediaDownload.FindByDownloadID(ctx, downloadID)
}

func (d *mediaDownload) GetAllFullNames(ctx context.Context) (map[string]struct{}, error) {
	return d.MediaDownload.GetAllFullNames(ctx, false)
}

func (d *mediaDownload) GetAllFullNamesWithDeleted(ctx context.Context) (map[string]struct{}, error) {
	return d.MediaDownload.GetAllFullNames(ctx, true)
}

func (uc *downloader) MediaDownload() MediaDownload {
	return &mediaDownload{
		MediaDownload: uc.download,
	}
}
