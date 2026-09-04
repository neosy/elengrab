package required

import (
	"context"
	"sync"
)

var (
	buildSearchIndexOnceSync sync.Once
)

func (m *migrations) buildSearchIndexOnce(ctx context.Context) (bool, error) {
	var (
		ok, executed bool
		err          error
	)

	buildSearchIndexOnceSync.Do(func() {
		err = m.buildSearchIndex(ctx)
		ok = err == nil
		executed = true
	})

	if executed {
		return ok, err
	}

	return true, nil
}

func (m *migrations) buildSearchIndex(ctx context.Context) error {
	return m.Usecases().SearchIndex.Build(ctx)
}
