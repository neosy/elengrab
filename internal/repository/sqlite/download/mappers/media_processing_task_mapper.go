package mappers

import (
	"database/sql"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	usql "github.com/neosy/elengrab/internal/pkg/dbutils/sql"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapMediaProcessingTaskToEntity(task *ddownload.MediaProcessingTask) (*edownload.MediaProcessingTask, error) {
	var jobID string
	if task.JobID != nil {
		jobID = task.JobID.String()
	}

	return &edownload.MediaProcessingTask{
		TaskID:     sql.NullString{String: task.TaskID.String(), Valid: true},
		DownloadID: sql.NullString{String: task.DownloadID.String(), Valid: true},
		Status:     sql.NullString{String: task.Status.String(), Valid: true},
		SourcePath: sql.NullString{String: task.SourcePath, Valid: true},
		WorkerID:   sql.NullInt64{Int64: int64(uptr.Deref(task.WorkerID)), Valid: task.WorkerID != nil},
		JobID:      sql.NullString{String: jobID, Valid: task.JobID != nil},
	}, nil
}

func (m *Mappers) MapMediaProcessingTaskEntityToDomain(
	eTask *edownload.MediaProcessingTask,
) (*ddownload.MediaProcessingTask, error) {
	status, err := dtypes.ParseMediaTaskStatus(usql.String(eTask.Status))
	if err != nil {
		return nil, err
	}

	taskID, err := uuid.Parse(usql.String(eTask.TaskID))
	if err != nil {
		return nil, err
	}

	downloadID, err := uuid.Parse(usql.String(eTask.DownloadID))
	if err != nil {
		return nil, err
	}

	var jobID *uuid.UUID
	if id := usql.String(eTask.JobID); id != "" {
		idd, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		jobID = uptr.Any(idd)
	}

	return &ddownload.MediaProcessingTask{
		TaskID:     taskID,
		DownloadID: downloadID,
		Status:     status,
		SourcePath: usql.String(eTask.SourcePath),
		WorkerID:   uptr.NonZero(uint64(usql.Int64(eTask.WorkerID))),
		JobID:      jobID,
		CreatedAt:  usql.Time(eTask.CreatedAt),
		UpdatedAt:  usql.Time(eTask.UpdatedAt),
	}, nil
}
