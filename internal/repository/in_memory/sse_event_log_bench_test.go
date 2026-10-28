package inmemory

import (
	"context"
	"testing"
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func BenchmarkPush(b *testing.B) {
	r := newSseEventLogRepository(time.Minute)
	ctx := context.Background()
	ev := &dtypes.SSEEvent{Name: "test", Data: []byte{}}

	// Fill the log up to the limit
	for range defaultLenSseEventLog {
		r.Push(ctx, "k", ev)
	}

	b.ReportAllocs()

	for b.Loop() {
		r.Push(ctx, "k", ev)
	}
}
