package itask

import (
	"context"
	"fmt"

	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type Task interface {
	DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error
}

// SelectTask returns the Task implementation for the given source type.
func SelectTask(
	sourceType dtypes.MediaDownloadSourceType,
	downloadTask Task,
	processingTask Task,
) (Task, error) {
	switch sourceType {
	case dtypes.MediaDownloadSourceTypeURL:
		return downloadTask, nil
	case dtypes.MediaDownloadSourceTypeUpload:
		return processingTask, nil
	default:
		return nil, fmt.Errorf("unsupported media download source type: %q", sourceType)
	}
}
