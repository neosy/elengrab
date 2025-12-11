package mappers

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
	usql "github.com/neosy/elengrab/pkg/utils/sql"
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
		TaskId:     sql.NullString{String: task.TaskId.String(), Valid: true},
		FileId:     sql.NullString{String: task.FileId.String(), Valid: true},
		Status:     sql.NullString{String: task.Status.String(), Valid: true},
		YoutubeUrl: sql.NullString{String: task.YoutubeUrl, Valid: true},
		Options:    sql.NullString{String: optionsJson, Valid: task.Options != nil},
		WorkerId:   sql.NullInt64{Int64: int64(uptr.Deref(task.WorkerId)), Valid: task.WorkerId != nil},
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
		jobID = uptr.Any(uuid.MustParse(id))
	}

	return &ddownload.DownloadTask{
		TaskId:     uuid.MustParse(usql.String(eTask.TaskId)),
		FileId:     uuid.MustParse(usql.String(eTask.FileId)),
		Status:     dtypes.DownloadTaskStatus(usql.String(eTask.Status)),
		YoutubeUrl: usql.String(eTask.YoutubeUrl),
		Options:    options,
		JobID:      jobID,
		WorkerId:   uptr.NonZero(uint(usql.Int64(eTask.WorkerId))),
		CreatedAt:  usql.Time(eTask.CreatedAt),
		UpdatedAt:  usql.Time(eTask.UpdatedAt),
	}, nil
}
