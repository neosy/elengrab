package mappers

import (
	"database/sql"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
	usql "github.com/neosy/elengrab/pkg/utils/sql"
)

func (m *Mappers) MapDownloadTaskDomainToEntity(task *ddownload.DownloadTask) *edownload.DownloadTask {
	return &edownload.DownloadTask{
		TaskId:   sql.NullString{String: task.TaskId.String(), Valid: true},
		FileId:   sql.NullString{String: task.FileId.String(), Valid: true},
		Status:   sql.NullString{String: task.Status.String(), Valid: true},
		WorkerId: sql.NullInt64{Int64: int64(uptr.Deref(task.WorkerId)), Valid: task.WorkerId != nil},
	}
}

func (m *Mappers) MapDownloadTaskEntityToDomain(eTask *edownload.DownloadTask) *ddownload.DownloadTask {
	return &ddownload.DownloadTask{
		TaskId:    uuid.MustParse(usql.String(eTask.TaskId)),
		FileId:    uuid.MustParse(usql.String(eTask.FileId)),
		Status:    dtypes.DownloadTaskStatus(usql.String(eTask.Status)),
		WorkerId:  uptr.NonZero(uint(usql.Int64(eTask.WorkerId))),
		CreatedAt: usql.Time(eTask.CreatedAt),
		UpdatedAt: usql.Time(eTask.UpdatedAt),
	}
}
