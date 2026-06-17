package downloader

import (
	"context"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

// GetDownloadInfo retrieves download information by download ID for a specific user.
func (uc *Downloader) GetDownloadInfo(
	ctx context.Context,
	authCtx dauth.UserContext,
	downloadID uuid.UUID,
) (*dto.GetMediaDownloadInfoResponse, error) {
	resp, err := uc.findActualDownloadInfo(ctx, nil, downloadID)
	if err != nil {
		uc.logger.Error("Failed get download info", "error", err)
		return nil, err
	}
	if resp == nil {
		return nil, apperrors.ErrDownloadNotFound
	}

	if !uc.authz.HasMediaViewAccess(authCtx, resp.MediaDownload) {
		return nil, ierrors.ErrAccessDenied
	}

	resp.HasWriteAccess = uc.HasWriteOperation(authCtx)

	return resp, nil
}

func (uc *Downloader) GetDownloadInfoUnrestricted(
	ctx context.Context,
	downloadID uuid.UUID,
) (*dto.GetMediaDownloadInfoResponse, error) {
	resp, err := uc.findActualDownloadInfo(ctx, nil, downloadID)
	if err != nil {
		uc.logger.Error("Failed get media download info", "error", err)
		return nil, err
	}
	if resp == nil {
		return nil, apperrors.ErrDownloadNotFound
	}
	return resp, nil
}

func (uc *Downloader) GetDownloadInfoForEdit(
	ctx context.Context,
	authCtx dauth.UserContext,
	downloadID uuid.UUID,
) (*dto.GetMediaDownloadInfoResponse, error) {
	resp, err := uc.GetDownloadInfo(ctx, authCtx, downloadID)
	if err != nil {
		return nil, err
	}

	err = uc.validateDownloadWriteAccess(authCtx, resp.MediaDownload)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// findActualDownloadInfo retrieves the actual download information based on user ID and download ID.
func (uc *Downloader) findActualDownloadInfo(
	ctx context.Context,
	userID *uuid.UUID,
	downloadID uuid.UUID,
) (*dto.GetMediaDownloadInfoResponse, error) {
	var download *ddownload.MediaDownload

	if downloadID == uuid.Nil {
		uc.logger.Warn("Id for the DownloadID field is not defined")
		return nil, apperrors.ErrDownloadIDIsNil
	}

	state, _ := uc.dlStateCache.FindByDownloadID(ctx, userID, downloadID)
	if state != nil && state.Download != nil {
		download = state.Download.Copy()
	}

	if download == nil {
		var err error
		download, err = uc.download.FindByDownloadID(ctx, userID, downloadID)
		if err != nil {
			return nil, err
		}
	}

	if download == nil {
		return nil, nil
	}

	return uc.findActualDownloadInfoByDownload(ctx, download)
}

// findActualDownloadInfoByDownload retrieves the actual download information based on the provided download.
func (uc *Downloader) findActualDownloadInfoByDownload(
	ctx context.Context,
	download *ddownload.MediaDownload,
) (*dto.GetMediaDownloadInfoResponse, error) {
	var (
		dlProgress *dservices.DownloaderProgress
	)

	if download == nil {
		return nil, nil
	}

	state, _ := uc.dlStateCache.FindByDownloadID(ctx, download.UserID, download.DownloadID)
	if state != nil && state.Download != nil {
		download = state.Download.Copy()
		dlProgress = state.Progress.Copy()
	}

	var avatarTitle string
	if download.ChannelID != nil && download.IsYouTube() {
		channel, _ := uc.ytChannel.FindByChannelID(ctx, *download.ChannelID)
		if channel != nil {
			avatarTitle = channel.ChannelTitle
		}
	}

	var hasSiteIcon bool
	siteLogo, _ := uc.siteIcon.FindBySiteURL(ctx, httpx.BaseURL(download.MediaURL))
	if siteLogo != nil {
		hasSiteIcon = true
		if avatarTitle == "" {
			avatarTitle = siteLogo.SiteTitle
		}
	}

	var isPortrait bool
	thumbnail, _ := uc.thumbnail.FindInfoByThumbID(ctx, download.MediaInfo.PreferredThumbnailID())
	if thumbnail != nil {
		isPortrait = thumbnail.IsPortrait()
	}

	return uc.mappers.MapDownloadDomainToDownloadInfoResponse(download, avatarTitle, dlProgress, hasSiteIcon, isPortrait), nil
}

func (uc *Downloader) getDownloadsInfo(
	ctx context.Context,
	authCtx dauth.UserContext,
	queryOptions dtypes.QueryOptions,
	filters map[string]any,
) ([]*dto.GetMediaDownloadInfoResponse, error) {
	var resps []*dto.GetMediaDownloadInfoResponse

	downloads, err := uc.download.GetBeforeTime(ctx, queryOptions, filters)
	if err != nil {
		uc.logger.Warn("Failed get downloads", "queryOptions", queryOptions, "error", err)
		return nil, err
	}

	resps = make([]*dto.GetMediaDownloadInfoResponse, 0, len(downloads))
	for _, download := range downloads {
		resp, err := uc.findActualDownloadInfoByDownload(ctx, download)
		if err != nil {
			continue
		}
		resp.HasWriteAccess = uc.HasWriteOperation(authCtx)
		resps = append(resps, resp)
	}

	return resps, nil
}

func (uc *Downloader) GetDownloadFilePath(ctx context.Context, downloadID uuid.UUID) (string, error) {
	download, err := uc.download.GetByDownloadID(ctx, nil, downloadID)
	if err != nil {
		uc.logger.Error("Failed find download", "error", err)
		return "", err
	}

	return uc.downloadsStorage.Path(download.FileFullName), nil
}

// GetDownloadFileName retrieves the display file name and extension
// for the given file ID.
//
// Returns:
//
//	filename - the human-readable name of the file
//	ext      - the file extension (without dot)
//	err      - an error if the record is not found or a query fails
func (uc *Downloader) GetDownloadFileName(
	ctx context.Context,
	userCtx dauth.UserContext,
	downloadID uuid.UUID,
) (string, string, error) {
	var accessByUserID *uuid.UUID
	if uc.authz.ShouldRestrictDownloads(userCtx.RoleIDs) {
		accessByUserID = &userCtx.UserID
	}

	download, err := uc.download.GetByDownloadID(ctx, accessByUserID, downloadID)
	if err != nil {
		uc.logger.Error("Failed find download", "error", err)
		return "", "", err
	}

	return download.SafeReadableFullName, download.Ext, nil
}
