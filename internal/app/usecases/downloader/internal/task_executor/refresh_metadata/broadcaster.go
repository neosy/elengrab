package refreshmetaexecutor

import (
	"context"

	"uuid"
)

type Broadcaster struct {
	MediaUpdated func(context.Context, uuid.UUID)
}
