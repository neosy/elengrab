package mappers

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

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

	return &edownload.MediaDownload{
		DownloadID:               download.DownloadID,
		UserID:                   download.UserID,
		Status:                   download.Status.String(),
		MediaURL:                 download.MediaURL,
		MediaTitle:               download.MediaTitle,
		MediaTitleLower:          strings.ToLower(download.MediaTitle),
		MediaDescription:         download.MediaDescription,
		ChannelID:                download.ChannelID,
		FileName:                 download.FileName,
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

func (m *Mappers) MapDownloadEntityToDomain(eDownload *edownload.MediaDownload, eTask *edownload.DownloadTask) (*ddownload.MediaDownload, error) {
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

	var task *ddownload.DownloadTask
	if eTask != nil && eTask.TaskID.String != "" {
		var err error
		task, err = m.MapDownloadTaskEntityToDomain(eTask)
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
		UserID:                   eDownload.UserID,
		Status:                   dtypes.MustParseMediaDownloadStatus(eDownload.Status),
		MediaURL:                 eDownload.MediaURL,
		MediaTitle:               eDownload.MediaTitle,
		MediaDescription:         eDownload.MediaDescription,
		ChannelID:                eDownload.ChannelID,
		FileName:                 eDownload.FileName,
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
		DownloadTask:             task,
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

		download, err := m.MapDownloadEntityToDomain(&eDownload, nil)
		if err != nil {
			return nil, err
		}

		downloads = append(downloads, download)
	}

	return downloads, nil
}

func (m *Mappers) MapRowsToDownloadsTask(rows *sql.Rows, fn func(*ddownload.MediaDownload) error) error {
	var (
		eDownload edownload.MediaDownload
		eTask     edownload.DownloadTask
	)

	for rows.Next() {
		err := rows.Scan(append(eDownload.FieldPointers(), eTask.FieldPointers()...)...)
		if err != nil {
			return err
		}

		download, err := m.MapDownloadEntityToDomain(&eDownload, &eTask)
		if err != nil {
			return err
		}

		err = fn(download)
		if err != nil {
			return err
		}
	}

	return nil
}
