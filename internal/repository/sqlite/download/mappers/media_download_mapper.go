package mappers

import (
	"database/sql"
	"encoding/json"
	"time"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapDownloadDomainToEntity(download *ddownload.MediaDownload) (*edownload.MediaDownload, error) {
	var mediaInfoJson string
	if download.MediaInfo != nil {
		data, err := json.MarshalIndent(download.MediaInfo, "", "  ")
		if err != nil {
			return nil, err
		}
		mediaInfoJson = string(data)
	}

	var downloadedAt *string
	if download.DownloadedAt != nil {
		downloadedAt = uptr.String(download.DownloadedAt.Format(dbutils.SQLiteDateTimeFormat))
	}

	var channelID *string
	if download.ChannelID != uuid.Nil() {
		channelID = new(download.ChannelID.String())
	}

	return &edownload.MediaDownload{
		DownloadID:               download.DownloadID,
		DownloadCode:             download.DownloadCode,
		SourceType:               download.SourceType.String(),
		UserID:                   download.UserID,
		Status:                   download.Status.String(),
		MediaURL:                 download.MediaURL,
		MediaTitleOriginal:       download.MediaTitleOriginal,
		MediaTitle:               download.MediaTitle,
		MediaDescriptionOriginal: download.MediaDescriptionOriginal,
		MediaDescription:         download.MediaDescription,
		ChannelID:                channelID,
		FileNameNoExt:            download.FileNameNoExt,
		Ext:                      download.Ext,
		FileFullName:             download.FileFullName,
		FileSize:                 download.FileSize,
		PartialHash:              download.PartialHash,
		SafeReadableFileFullName: download.SafeReadableFileFullName,
		MediaInfo:                &mediaInfoJson,
		ErrorMessage:             download.ErrorMessage,
		Visibility:               download.Visibility.String(),
		DownloadedAt:             downloadedAt,
	}, nil
}

func (m *Mappers) MapDownloadEntityToDomain(
	eDownload *edownload.MediaDownload,
	eDownloadTask *edownload.DownloadTask,
	eProcessingTask *edownload.MediaProcessingTask,
) (*ddownload.MediaDownload, error) {
	sourceType, err := dtypes.ParseMediaDownloadSourceType(eDownload.SourceType)
	if err != nil {
		return nil, err
	}

	var mediaInfo *dtypes.MediaInfo
	if eDownload.MediaInfo != nil {
		mediaInfoJson := *eDownload.MediaInfo
		if mediaInfoJson != "" {
			mediaInfo = &dtypes.MediaInfo{}
			err := json.Unmarshal([]byte(mediaInfoJson), mediaInfo)
			if err != nil {
				return nil, err
			}
			if mediaInfo.AudioInfo != nil {
				mediaInfo.AudioInfo.Codec = dtypes.MustParseAudioCodec(string(mediaInfo.AudioInfo.Codec))
			}
			if mediaInfo.VideoInfo != nil {
				mediaInfo.VideoInfo.Codec = dtypes.MustParseVideoCodec(string(mediaInfo.VideoInfo.Codec))
			}
		}
	}

	var channelID uuid.UUID
	if eDownload.ChannelID != nil && *eDownload.ChannelID != "" {
		id, err := uuid.Parse(*eDownload.ChannelID)
		if err == nil {
			channelID = id
		}
	}

	var downloadTask *ddownload.DownloadTask
	if eDownloadTask != nil && eDownloadTask.TaskID.String != "" {
		var err error
		downloadTask, err = m.MapDownloadTaskEntityToDomain(eDownloadTask)
		if err != nil {
			return nil, err
		}
	}

	var processingTask *ddownload.MediaProcessingTask
	if eProcessingTask != nil && eProcessingTask.TaskID.String != "" {
		var err error
		processingTask, err = m.MapMediaProcessingTaskEntityToDomain(eProcessingTask)
		if err != nil {
			return nil, err
		}
	}

	var downloadedAt *time.Time
	if eDownload.DownloadedAt != nil {
		t, err := time.Parse(time.RFC3339, *eDownload.DownloadedAt)
		if err != nil {
			return nil, err
		}
		downloadedAt = &t
	}

	visibility, err := dtypes.ParseMediaVisibility(eDownload.Visibility)
	if err != nil {
		return nil, err
	}

	return &ddownload.MediaDownload{
		DownloadID:               eDownload.DownloadID,
		DownloadCode:             eDownload.DownloadCode,
		SourceType:               sourceType,
		UserID:                   eDownload.UserID,
		Status:                   dtypes.MustParseMediaDownloadStatus(eDownload.Status),
		MediaURL:                 eDownload.MediaURL,
		MediaTitleOriginal:       eDownload.MediaTitleOriginal,
		MediaTitle:               eDownload.MediaTitle,
		MediaDescriptionOriginal: eDownload.MediaDescriptionOriginal,
		MediaDescription:         eDownload.MediaDescription,
		ChannelID:                channelID,
		FileNameNoExt:            eDownload.FileNameNoExt,
		Ext:                      eDownload.Ext,
		FileFullName:             eDownload.FileFullName,
		FileSize:                 eDownload.FileSize,
		PartialHash:              eDownload.PartialHash,
		SafeReadableFileFullName: eDownload.SafeReadableFileFullName,
		MediaInfo:                mediaInfo,
		ErrorMessage:             eDownload.ErrorMessage,
		Visibility:               visibility,
		DownloadedAt:             downloadedAt,
		CreatedAt:                eDownload.CreatedAt,
		UpdatedAt:                eDownload.UpdatedAt,
		DeletedAt:                eDownload.DeletedAt,
		DownloadTask:             downloadTask,
		ProcessingTask:           processingTask,
	}, nil
}

func (m *Mappers) MapRowsToDownloads(rows *sql.Rows) ([]*ddownload.MediaDownload, error) {
	var (
		eDownload edownload.MediaDownload
		downloads []*ddownload.MediaDownload
	)

	for rows.Next() {
		err := rows.Scan(eDownload.FieldPointers()...)
		if err != nil {
			return nil, err
		}

		download, err := m.MapDownloadEntityToDomain(&eDownload, nil, nil)
		if err != nil {
			return nil, err
		}

		downloads = append(downloads, download)
	}

	return downloads, nil
}
