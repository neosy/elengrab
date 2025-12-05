package mappers

import (
	"database/sql"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapFileDomainToEntity(file *ddownload.File) *edownload.File {
	return &edownload.File{
		FileId:               file.FileId,
		Status:               file.Status.String(),
		YoutubeUrl:           file.YoutubeUrl,
		YoutubeTitle:         file.YoutubeTitle,
		FileName:             file.FileName,
		Ext:                  file.Ext,
		FullName:             file.FullName,
		FileSize:             file.FileSize,
		PartialHash:          file.PartialHash,
		SafeReadableFullName: file.SafeReadableFullName,
		ErrorMessage:         file.ErrorMessage,
	}
}

func (m *Mappers) MapFileEntityToDomain(eFile *edownload.File, eTask *edownload.DownloadTask) (*ddownload.File, error) {
	var task *ddownload.DownloadTask
	if eTask != nil && eTask.TaskId.String != "" {
		var err error
		task, err = m.MapDownloadTaskEntityToDomain(eTask)
		if err != nil {
			return nil, err
		}
	}

	return &ddownload.File{
		FileId:               eFile.FileId,
		Status:               dtypes.FileStatus(eFile.Status),
		YoutubeUrl:           eFile.YoutubeUrl,
		YoutubeTitle:         eFile.YoutubeTitle,
		FileName:             eFile.FileName,
		Ext:                  eFile.Ext,
		FullName:             eFile.FullName,
		FileSize:             eFile.FileSize,
		PartialHash:          eFile.PartialHash,
		SafeReadableFullName: eFile.SafeReadableFullName,
		ErrorMessage:         eFile.ErrorMessage,
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

func (m *Mappers) MapRowsToFilesTask(rows *sql.Rows) ([]*ddownload.File, error) {
	var (
		eFile edownload.File
		eTask edownload.DownloadTask
		files []*ddownload.File
	)

	for rows.Next() {
		err := rows.Scan(append(eFile.FieldPointers(), eTask.FieldPointers()...)...)
		if err != nil {
			return nil, err
		}

		file, err := m.MapFileEntityToDomain(&eFile, &eTask)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}
