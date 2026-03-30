package downloader

import (
	"context"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
)

// GetFileInfo retrieves file information by file ID for a specific user.
func (uc *Downloader) GetFileInfo(
	ctx context.Context,
	userCtx dauth.UserContext,
	fileID uuid.UUID,
) (*dto.GetFileInfoResponse, error) {
	var accessByUserID *uuid.UUID
	if uc.authz.RestrictFilesByUser(userCtx.Roles) {
		accessByUserID = &userCtx.UserID
	}

	resp, err := uc.findActualFileInfo(ctx, accessByUserID, fileID)
	if err != nil {
		uc.logger.Error("Failed get file info", "error", err)
		return nil, err
	}
	if resp == nil {
		return nil, errorx.New("file not found", exceptionx.NOT_FOUND)
	}
	return resp, nil
}

// findActualFileInfo retrieves the actual file information based on user ID and file ID.
func (uc *Downloader) findActualFileInfo(
	ctx context.Context,
	userID *uuid.UUID,
	fileID uuid.UUID,
) (*dto.GetFileInfoResponse, error) {
	var file *ddownload.File

	if fileID == uuid.Nil {
		uc.logger.Warn("Id for the FileID field is not defined")
		return nil, errorx.New("fileId not specified", exceptionx.ERROR)
	}

	state, _ := uc.dlStateCache.FindByFileID(ctx, userID, fileID)
	if state != nil && state.File != nil {
		file = state.File.Copy()
	}

	if file == nil {
		var err error
		file, err = uc.file.FindByFileID(ctx, userID, fileID)
		if err != nil {
			return nil, err
		}
	}

	if file == nil {
		return nil, nil
	}

	return uc.findActualFileInfoByFile(ctx, file)
}

// findActualFileInfoByFile retrieves the actual file information based on the provided file.
func (uc *Downloader) findActualFileInfoByFile(
	ctx context.Context,
	file *ddownload.File,
) (*dto.GetFileInfoResponse, error) {
	var (
		dlProgress *ddownload.DownloadProgress
	)

	if file == nil {
		return nil, nil
	}

	state, _ := uc.dlStateCache.FindByFileID(ctx, file.UserID, file.FileID)
	if state != nil && state.File != nil {
		file = state.File.Copy()
		dlProgress = state.Progress.Copy()
	}

	var avatarTitle string
	if file.YoutubeChannelID != nil {
		channel, _ := uc.ytChannel.FindByChannelID(ctx, *file.YoutubeChannelID)
		if channel != nil {
			avatarTitle = channel.ChannelTitle
		}
	}

	var hasSiteIcon bool
	siteLogo, _ := uc.siteIcon.FindBySiteURL(ctx, httpx.BaseURL(file.MediaUrl))
	if siteLogo != nil {
		hasSiteIcon = true
		if avatarTitle == "" {
			avatarTitle = siteLogo.SiteTitle
		}
	}

	return uc.mappers.MapFileDomainToFileInfoResponse(file, avatarTitle, dlProgress, uc.downloadsDir, hasSiteIcon), nil
}

// LoadHistory retrieves the download history for a user.
func (uc *Downloader) LoadHistory(
	ctx context.Context,
	userCtx dauth.UserContext,
	before time.Time,
	limit uint64,
	filterByTitle string,
) ([]*dto.GetFileInfoResponse, error) {
	filters := make(map[string]any)
	if uc.authz.RestrictFilesByUser(userCtx.Roles) {
		filters["userID"] = userCtx.UserID
	}

	if filterByTitle != "" {
		filters["title"] = filterByTitle
	}

	return uc.getFilesInfo(ctx, before, limit, filters)
}

func (uc *Downloader) getFilesInfo(
	ctx context.Context,
	before time.Time,
	limit uint64,
	filters map[string]any,
) ([]*dto.GetFileInfoResponse, error) {
	var resps []*dto.GetFileInfoResponse

	files, err := uc.file.GetBeforeTime(ctx, before, limit, filters)
	if err != nil {
		uc.logger.Warn("Failed get files", "before", before, "limit", limit, "error", err)
		return nil, err
	}

	resps = make([]*dto.GetFileInfoResponse, 0, len(files))
	for _, file := range files {
		resp, err := uc.findActualFileInfoByFile(ctx, file)
		if err != nil {
			continue
		}
		resps = append(resps, resp)
	}

	return resps, nil
}

func (uc *Downloader) GetFilePath(ctx context.Context, fileId uuid.UUID) (string, error) {
	file, err := uc.file.GetByFileID(ctx, nil, fileId)
	if err != nil {
		uc.logger.Error("Failed find file", "error", err)
		return "", err
	}

	return filepath.Join(uc.downloadsDir, file.FullName), nil
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
	fileId uuid.UUID,
) (string, string, error) {
	var accessByUserID *uuid.UUID
	if uc.authz.RestrictFilesByUser(userCtx.Roles) {
		accessByUserID = &userCtx.UserID
	}

	file, err := uc.file.GetByFileID(ctx, accessByUserID, fileId)
	if err != nil {
		uc.logger.Error("Failed find file", "error", err)
		return "", "", err
	}

	return file.SafeReadableFullName, file.Ext, nil
}
