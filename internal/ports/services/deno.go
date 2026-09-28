package pservices

import (
	"context"
)

type Deno interface {
	GetVersion(ctx context.Context) (string, error)
}
