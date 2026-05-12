package mappers

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	usql "github.com/neosy/elengrab/internal/pkg/dbutils/sql"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapDownloadTaskDomainToEntity(task *ddownload.DownloadTask) (*edownload.DownloadTask, error) {
	var optionsJson string
	if task.Options != nil {
		data, err := json.MarshalIndent(task.Options, "", "  ")
		if err != nil {
			return nil, err
		}
		optionsJson = string(data)
	}

	var jobID string
	if task.JobID != nil {
		jobID = task.JobID.String()
	}

	return &edownload.DownloadTask{
		TaskID:     sql.NullString{String: task.TaskID.String(), Valid: true},
		DownloadID: sql.NullString{String: task.DownloadID.String(), Valid: true},
		Status:     sql.NullString{String: task.Status.String(), Valid: true},
		MediaUrl:   sql.NullString{String: task.MediaUrl, Valid: true},
		Options:    sql.NullString{String: optionsJson, Valid: task.Options != nil},
		WorkerID:   sql.NullInt64{Int64: int64(uptr.Deref(task.WorkerID)), Valid: task.WorkerID != nil},
		JobID:      sql.NullString{String: jobID, Valid: task.JobID != nil},
	}, nil
}

func (m *Mappers) MapDownloadTaskEntityToDomain(eTask *edownload.DownloadTask) (*ddownload.DownloadTask, error) {
	var options *ddownload.DownloadOptions

	optionsJson := usql.String(eTask.Options)
	if optionsJson != "" {
		options = &ddownload.DownloadOptions{}
		err := json.Unmarshal([]byte(optionsJson), options)
		if err != nil {
			return nil, err
		}
	}

	var jobID *uuid.UUID
	if id := usql.String(eTask.JobID); id != "" {
		idd, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		jobID = uptr.Any(idd)
	}

	taskID, err := uuid.Parse(usql.String(eTask.TaskID))
	if err != nil {
		return nil, err
	}

	fileID, err := uuid.Parse(usql.String(eTask.DownloadID))
	if err != nil {
		return nil, err
	}

	return &ddownload.DownloadTask{
		TaskID:     taskID,
		DownloadID: fileID,
		Status:     dtypes.MustParseDownloadTaskStatus(usql.String(eTask.Status)),
		MediaUrl:   usql.String(eTask.MediaUrl),
		Options:    options,
		JobID:      jobID,
		WorkerID:   uptr.NonZero(uint64(usql.Int64(eTask.WorkerID))),
		CreatedAt:  usql.Time(eTask.CreatedAt),
		UpdatedAt:  usql.Time(eTask.UpdatedAt),
	}, nil
}
