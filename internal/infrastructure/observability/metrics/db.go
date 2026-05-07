package metrics

func SetTableRows(dbName, table string, count int64) {
	dbTableRows.WithLabelValues(dbName, table).Set(float64(count))
}
