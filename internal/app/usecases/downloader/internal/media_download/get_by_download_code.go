package mediadownload

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaDownload) FindByDownloadCodeNoCache(
	ctx context.Context,
	downloadCode string,
) (*ddownload.MediaDownload, error) {
	if downloadCode == "" {
		return nil, nil
	}

	download, err := uc.downloadRepo().FindByDownloadCode(ctx, downloadCode)
	if err != nil {
		uc.logger.Warn("Failed to find record", "error", err)
		return nil, err
	}

	return download, err
}

func (uc *MediaDownload) ExistsByDownloadCodeNoCache(
	ctx context.Context,
	downloadCode string,
) (bool, error) {
	if downloadCode == "" {
		return false, nil
	}

	exists, err := uc.downloadRepo().ExistsByDownloadCode(ctx, downloadCode)
	if err != nil {
		uc.logger.Warn("Failed to find record", "error", err)
		return false, err
	}

	return exists, err
}

func (uc *MediaDownload) GetByDownloadCodeNoCache(
	ctx context.Context,
	downloadCode string,
) (*ddownload.MediaDownload, error) {
	if downloadCode == "" {
		uc.logger.Warn("MediaDownload not found", "downloadCode", downloadCode)
		return nil, nil
	}

	download, err := uc.FindByDownloadCodeNoCache(ctx, downloadCode)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if download == nil {
		uc.logger.Warn("MediaDownload not found", "downloadCode", downloadCode)
		return nil, errorx.New("download not found", exceptionx.NOT_FOUND)
	}

	return download, nil
}

func (uc *MediaDownload) FindByDownloadCode(
	ctx context.Context,
	downloadCode string,
) (*ddownload.MediaDownload, error) {
	if downloadCode == "" {
		return nil, nil
	}

	mediaDownload, cacheStatus, _ := uc.downloadCacheRep.FindByDownloadCode(ctx, downloadCode)
	if mediaDownload != nil {
		return mediaDownload, nil
	}
	if cacheStatus == memsimple.CacheStatusNegativeHit {
		return nil, nil
	}

	mediaDownload, err := uc.FindByDownloadCodeNoCache(ctx, downloadCode)
	if err != nil {
		return nil, err
	}

	if mediaDownload == nil {
		uc.downloadCacheRep.SaveNegativeByCode(ctx, downloadCode)
		return nil, nil
	}

	uc.downloadCacheRep.Save(ctx, mediaDownload)

	return mediaDownload, nil
}

func (uc *MediaDownload) GetByDownloadCode(
	ctx context.Context,
	downloadCode string,
) (*ddownload.MediaDownload, error) {
	if downloadCode == "" {
		uc.logger.Warn("MediaDownload not found", "downloadCode", downloadCode)
		return nil, nil
	}

	download, err := uc.FindByDownloadCode(ctx, downloadCode)
	if err != nil {
		return nil, err
	}

	if download == nil {
		uc.logger.Warn("MediaDownload not found", "downloadCode", downloadCode)
		return nil, errorx.New("download not found", exceptionx.NOT_FOUND)
	}

	return download, nil
}
