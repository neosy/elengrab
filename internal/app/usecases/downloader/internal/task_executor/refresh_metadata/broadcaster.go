package refreshmetaexecutor

import (
	"context"

	"github.com/google/uuid"
)

type Broadcaster struct {
	MediaUpdated func(context.Context, uuid.UUID)
}
