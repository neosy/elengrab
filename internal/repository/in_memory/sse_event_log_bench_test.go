package inmemory

import (
	"context"
	"strconv"
	"testing"
	"time"

	dsse "github.com/neosy/elengrab/internal/domain/types/sse"
)

func BenchmarkPush(b *testing.B) {
	r := newSseEventLogRepository(time.Minute)
	ctx := context.Background()
	ev := &dsse.Event{Name: "test", Data: []byte{}}

	// Fill the log up to the limit
	for i := range defaultLenSseEventLog {
		ev.Data = []byte(strconv.Itoa(i))
		r.Push(ctx, "k", ev)
	}

	datas := make([][]byte, defaultLenSseEventLog*2)

	for i := range defaultLenSseEventLog * 2 {
		datas[i] = []byte(strconv.Itoa(i))
	}

	b.ReportAllocs()

	var i int
	for b.Loop() {
		if i == defaultLenSseEventLog*2-1 {
			i = 0
		}

		ev.Data = datas[i]
		r.Push(ctx, "k", ev)
	}
}
