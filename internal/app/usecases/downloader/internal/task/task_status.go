package itask

import (
	"context"
	"fmt"

	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type TaskStatus interface {
	New(ctx context.Context, taskId uuid.UUID) error
	Penging(ctx context.Context, taskId uuid.UUID, jobID uuid.UUID) error
	Working(ctx context.Context, taskId uuid.UUID, workerId uint64) error
	Failed(ctx context.Context, taskId uuid.UUID) error
}

// TaskStatusBySourceType returns the TaskStatus implementation for the given source type.
func TaskStatusBySourceType(
	sourceType dtypes.MediaDownloadSourceType,
	downloadTaskStatus TaskStatus,
	processingTaskStatus TaskStatus,
) (TaskStatus, error) {
	switch sourceType {
	case dtypes.MediaDownloadSourceTypeURL:
		return downloadTaskStatus, nil
	case dtypes.MediaDownloadSourceTypeUpload:
		return processingTaskStatus, nil
	default:
		return nil, fmt.Errorf("unsupported media download source type: %q", sourceType)
	}
}
