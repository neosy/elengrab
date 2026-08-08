package downloader

import (
	"context"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

// GetDownloadInfo retrieves download information by download ID for a specific user.
func (uc *Downloader) GetDownloadInfo(
	ctx context.Context,
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
) (*dto.MediaDownloadInfo, error) {
	resp, err := uc.findActualDownloadInfo(ctx, downloadID, withAuth(authCtx))
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
) (*dto.MediaDownloadInfo, error) {
	resp, err := uc.findActualDownloadInfo(ctx, downloadID)
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
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
) (*dto.MediaDownloadInfo, error) {
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
	downloadID uuid.UUID,
	opts ...callOption,
) (*dto.MediaDownloadInfo, error) {
	var download *ddownload.MediaDownload

	if downloadID == uuid.Nil {
		uc.logger.Warn("Id for the DownloadID field is not defined")
		return nil, apperrors.ErrDownloadIDIsNil
	}

	state, _ := uc.dlStateCache.FindByDownloadID(ctx, downloadID)
	if state != nil && state.Download != nil {
		download = state.Download.Copy()
	}

	if download == nil {
		var err error
		download, err = uc.download.FindByDownloadIDNoCache(ctx, downloadID)
		if err != nil {
			return nil, err
		}
	}

	if download == nil {
		return nil, nil
	}

	return uc.findActualDownloadInfoByDownload(ctx, download, opts...)
}

// findActualDownloadInfoByDownload retrieves the actual download information based on the provided download.
func (uc *Downloader) findActualDownloadInfoByDownload(
	ctx context.Context,
	download *ddownload.MediaDownload,
	opts ...callOption,
) (*dto.MediaDownloadInfo, error) {
	var (
		dlProgress *dservices.DownloaderProgress
	)

	if download == nil {
		return nil, nil
	}

	state, _ := uc.dlStateCache.FindByDownloadID(ctx, download.DownloadID)
	if state != nil && state.Download != nil {
		download = state.Download.Copy()
		dlProgress = state.Progress.Copy()
	}

	var login string
	if download.UserID != nil {
		user, _ := uc.authSrv.FindByUserID(ctx, *download.UserID)
		if user != nil {
			login = user.Login.String()
		}
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
	thumbnail, _ := uc.thumbnail.FindByThumbID(ctx, download.MediaInfo.PreferredThumbnailID())
	if thumbnail != nil {
		isPortrait = thumbnail.IsPortrait()
	}

	viewCount, _ := uc.mediaWatch.GetViews(ctx, download.DownloadID)

	options := buildCallOptions(opts...)

	var authCtx dauth.AuthContext
	if options.authCtx != nil {
		authCtx = *options.authCtx
	}

	var lastWatchPosition time.Duration
	if options.authCtx != nil {
		userID := options.authCtx.UserID
		lastWatchPosition, _ = uc.mediaWatch.GetLastUserWatchPosition(ctx, download.DownloadID, userID, options.authCtx.AnonSessionID)
	}

	var watched bool
	if authCtx.UserID != uuid.Nil {
		watched, _ = uc.mediaWatch.HasUserWatched(ctx, download.DownloadID, authCtx.UserID)
	}

	mappingData := &dto.MediaDownloadInfoMappingData{
		UserLogin:   login,
		AvatarTitle: avatarTitle,

		ViewCount: viewCount,

		UserLastWatchPosition: lastWatchPosition,
		UserWatched:           watched,

		HasSiteIcon:         hasSiteIcon,
		ThumbnailIsPortrait: isPortrait,

		Progress: dlProgress,
	}

	return uc.mappers.MapDownloadDomainToDownloadInfoResponse(download, mappingData), nil
}

func (uc *Downloader) GetDownloadFilePath(ctx context.Context, downloadID uuid.UUID) (string, error) {
	download, err := uc.download.GetByDownloadIDNoCache(ctx, downloadID)
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
//	err      - an error if the record is not found or a query fails
func (uc *Downloader) GetDownloadFileName(
	ctx context.Context,
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
) (string, error) {
	download, err := uc.download.GetByDownloadIDNoCache(ctx, downloadID)
	if err != nil {
		uc.logger.Error("Failed find download", "error", err)
		return "", err
	}

	if !uc.authz.HasMediaViewAccess(authCtx, download) {
		return "", ierrors.ErrAccessDenied
	}

	return download.SafeReadableFileFullName, nil
}
