package downloader

import (
	"context"
	"time"

	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/broadcaster"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type downloaderBase interface {
	AppMode() dtypes.AppMode
	DemoMode() bool
}

type Downloader interface {
	downloaderBase

	Start(ctx context.Context)
}

type DownloaderMaintenance interface {
	downloaderBase

	UpdateHash(ctx context.Context) error
	DeleteDuplicates(ctx context.Context) error
	DeleteMissingMediaDownloads(ctx context.Context, enableMoveUnmatchedFiles bool) error
	DeleteFailedMediaDownloads(ctx context.Context) error
	ResetStuckJobs(ctx context.Context) error
	ResolveInconsistencies(ctx context.Context, orphanGracePeriod time.Duration) error
}

type DownloaderTask interface {
	downloaderBase

	ExecuteDownloadTask(
		ctx context.Context,
		workerID uint64,
		task *ddownload.DownloadTask,
	) error

	ExecuteRefreshMetadataTask(
		ctx context.Context,
		workerID uint64,
		task *ddownload.RefreshMetadataTask,
	) error

	UpdateSystemDiskInfo()
	UpdateSystemInfo()
}

type DownloaderAPI interface {
	downloaderBase

	MigrateGuestData(ctx context.Context, guestID, userID uuid.UUID) error
	Broadcaster() *broadcaster.Broadcaster
	NotifyDownloadUpdated(ctx context.Context, downloadID uuid.UUID)
	NotifyDownloadChanged(ctx context.Context, req dto.MediaDownloadChanged)
	FindChannelInfo(ctx context.Context, externalID string, platform string) (*dmedia.Channel, error)
	GetChannelInfo(ctx context.Context, externalID string, platform string) (*dmedia.Channel, error)
	GetChannelByID(ctx context.Context, channelID uuid.UUID) (*dmedia.Channel, error)
	CheckDownloadVisibilityAccess(ctx context.Context, authCtx dauth.AuthContext, downloadID uuid.UUID) error
	GetDownloadInfo(ctx context.Context, authCtx dauth.AuthContext, downloadID uuid.UUID) (*dto.MediaDownloadInfo, error)
	GetDownloadInfoUnrestricted(ctx context.Context, downloadID uuid.UUID) (*dto.MediaDownloadInfo, error)
	GetDownloadInfoForEdit(
		ctx context.Context,
		authCtx dauth.AuthContext,
		downloadID uuid.UUID,
	) (*dto.MediaDownloadInfo, error)
	ListDownloadInfo(
		ctx context.Context,
		authCtx dauth.AuthContext,
		query dto.MediaDownloadQuery,
	) ([]*dto.MediaDownloadInfo, error)
	PatchMediaDownload(ctx context.Context, authCtx dauth.AuthContext, req dto.PatchMediaDownloadRequest) error
	DeleteDownload(ctx context.Context, authCtx dauth.AuthContext, downloadID uuid.UUID) error
	HasWriteAccess(authCtx dauth.AuthContext) bool
	CanCreateMediaDownload(authCtx dauth.AuthContext) bool
	GetLastWatchPosition(ctx context.Context, authCtx dauth.AuthContext, downloadID uuid.UUID) (time.Duration, error)
	ScheduleDownload(
		ctx context.Context,
		authCtx dauth.AuthContext,
		url string,
		options *ddownload.DownloadOptions,
	) (*dto.ScheduleDownloadResponse, error)
	RetryDownload(
		ctx context.Context,
		authCtx dauth.AuthContext,
		downloadID uuid.UUID,
	) (*dto.MediaDownloadInfo, error)
	ScheduleUpload(
		ctx context.Context,
		authCtx dauth.AuthContext,
		upload pservices.UploadedFile,
		req dto.ScheduleUploadRequest,
	) (*dto.ScheduleUploadResponse, error)
	ScheduleRefreshMetadata(
		ctx context.Context,
		authCtx dauth.AuthContext,
		downloadID uuid.UUID,
	) error
	GetDownloadImage(
		ctx context.Context,
		userCtx dauth.AuthContext,
		downloadID uuid.UUID,
		sources []dtypes.ImageSource,
	) (*dtypes.ImageData, error)
	GetChannelImage(ctx context.Context, channel *dmedia.Channel) (*dtypes.ImageData, error)
	GetSiteImage(ctx context.Context, url string) (*dtypes.ImageData, error)
	TrackMediaWatchEvent(ctx context.Context, authCtx dauth.AuthContext, req dto.TrackMediaWatchEventRequest) error
	SystemInfo() dto.SystemInfoResponse
}

type InternalDownloader interface {
	downloaderBase

	// Usecases
	MediaDownload() MediaDownload
	MediaWatch() MediaWatch
	SearchIndex() SearchIndex
	Channel() Channel
}
