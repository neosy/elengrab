package dlexecutor

import (
	"context"

	"github.com/google/uuid"
)

type Broadcaster struct {
	MediaUpdated            func(context.Context, uuid.UUID)
	DownloadProgressUpdated func(context.Context, uuid.UUID)
}
