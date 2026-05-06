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

func (m *Mappers) MapFileDomainToEntity(file *ddownload.File) (*edownload.File, error) {
	var mediaInfoJson string
	if file.MediaInfo != nil {
		data, err := json.MarshalIndent(file.MediaInfo, "", "  ")
		if err != nil {
			return nil, err
		}
		mediaInfoJson = string(data)
	}

	var downloadedAt *string
	if file.DownloadedAt != nil {
		downloadedAt = uptr.String(file.DownloadedAt.Format(dbutils.SQLiteDateTimeFormat))
	}

	return &edownload.File{
		FileID:               file.FileID,
		UserID:               file.UserID,
		Status:               file.Status.String(),
		MediaUrl:             file.MediaUrl,
		MediaTitle:           file.MediaTitle,
		MediaTitleLower:      strings.ToLower(file.MediaTitle),
		ChannelID:            file.ChannelID,
		FileName:             file.FileName,
		Ext:                  file.Ext,
		FullName:             file.FullName,
		FileSize:             file.FileSize,
		PartialHash:          file.PartialHash,
		SafeReadableFullName: file.SafeReadableFullName,
		MediaInfo:            &mediaInfoJson,
		ErrorMessage:         file.ErrorMessage,
		DownloadedAt:         downloadedAt,
	}, nil
}

func (m *Mappers) MapFileEntityToDomain(eFile *edownload.File, eTask *edownload.DownloadTask) (*ddownload.File, error) {
	var mediaInfo *dtypes.MediaInfo
	if eFile.MediaInfo != nil {
		mediaInfoJson := *eFile.MediaInfo
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
	if eFile.DownloadedAt != nil {
		t, err := time.Parse(time.RFC3339, *eFile.DownloadedAt)
		if err != nil {
			return nil, err
		}
		downloadedAt = &t
	}

	return &ddownload.File{
		FileID:               eFile.FileID,
		UserID:               eFile.UserID,
		Status:               dtypes.MustParseFileStatus(eFile.Status),
		MediaUrl:             eFile.MediaUrl,
		MediaTitle:           eFile.MediaTitle,
		ChannelID:            eFile.ChannelID,
		FileName:             eFile.FileName,
		Ext:                  eFile.Ext,
		FullName:             eFile.FullName,
		FileSize:             eFile.FileSize,
		PartialHash:          eFile.PartialHash,
		SafeReadableFullName: eFile.SafeReadableFullName,
		MediaInfo:            mediaInfo,
		ErrorMessage:         eFile.ErrorMessage,
		DownloadedAt:         downloadedAt,
		CreatedAt:            eFile.CreatedAt,
		UpdatedAt:            eFile.UpdatedAt,
		DeletedAt:            eFile.DeletedAt,
		DownloadTask:         task,
	}, nil
}

func (m *Mappers) MapRowsToFiles(rows *sql.Rows) ([]*ddownload.File, error) {
	var (
		eFile edownload.File
		files []*ddownload.File
	)

	for rows.Next() {
		err := rows.Scan(eFile.FieldPointers()...)
		if err != nil {
			return nil, err
		}

		file, err := m.MapFileEntityToDomain(&eFile, nil)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

func (m *Mappers) MapRowsToFilesTask(rows *sql.Rows, fn func(*ddownload.File) error) error {
	var (
		eFile edownload.File
		eTask edownload.DownloadTask
	)

	for rows.Next() {
		err := rows.Scan(append(eFile.FieldPointers(), eTask.FieldPointers()...)...)
		if err != nil {
			return err
		}

		file, err := m.MapFileEntityToDomain(&eFile, &eTask)
		if err != nil {
			return err
		}

		err = fn(file)
		if err != nil {
			return err
		}
	}

	return nil
}
