package sqliterep

import "github.com/neosy/elengrab/internal/infrastructure/observability/metrics"

func (r *Repositories) UpdateMetrics() error {
	var count int64

	for dbName, entry := range r.dbRegistry.EntriesByName() {
		for _, name := range entry.Schema().TableNames() {
			err := entry.DB().QueryRow("SELECT COUNT(*) FROM " + name).Scan(&count)
			if err != nil {
				return err
			}
			metrics.SetTableRows(dbName, name, count)
		}
	}

	return nil
}
