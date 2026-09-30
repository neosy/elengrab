package required

import (
	"context"
	"sync"
)

var (
	buildSearchIndexOnceSync sync.Once
)

func (m *migrations) buildSearchIndexOnce(ctx context.Context) error {
	var err error

	buildSearchIndexOnceSync.Do(func() {
		err = m.buildSearchIndex(ctx)
	})

	return err
}

func (m *migrations) buildSearchIndex(ctx context.Context) error {
	return m.Usecases().SearchIndex.Build(ctx)
}
