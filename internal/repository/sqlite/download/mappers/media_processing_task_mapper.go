package mappers

import (
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapMediaProcessingTaskToEntity(task *ddownload.MediaProcessingTask) (*edownload.MediaProcessingTask, error) {
	var workerID *int
	if task.WorkerID != nil {
		workerID = new(int(*task.WorkerID))
	}

	return &edownload.MediaProcessingTask{
		TaskID:     task.TaskID,
		DownloadID: task.DownloadID,
		Status:     task.Status.String(),
		SourcePath: task.SourcePath,
		WorkerID:   workerID,
		JobID:      task.JobID,
	}, nil
}

func (m *Mappers) MapMediaProcessingTaskEntityToDomain(
	etask *edownload.MediaProcessingTask,
) (*ddownload.MediaProcessingTask, error) {
	status, err := dtypes.ParseMediaTaskStatus(etask.Status)
	if err != nil {
		return nil, err
	}

	var workerID *uint64
	if etask.WorkerID != nil {
		workerID = new(uint64(*etask.WorkerID))
	}

	return &ddownload.MediaProcessingTask{
		TaskID:     etask.TaskID,
		DownloadID: etask.DownloadID,
		Status:     status,
		SourcePath: etask.SourcePath,
		WorkerID:   workerID,
		JobID:      etask.JobID,
		CreatedAt:  etask.CreatedAt,
		UpdatedAt:  etask.UpdatedAt,
	}, nil
}
