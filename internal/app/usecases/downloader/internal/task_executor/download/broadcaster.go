package dlexecutor

import (
	"context"

	"uuid"
)

type Broadcaster struct {
	MediaUpdated            func(context.Context, uuid.UUID)
	DownloadProgressUpdated func(context.Context, uuid.UUID)
}
