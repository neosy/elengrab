package pworkers

type DBMMetricsRunner interface {
	UpdateMetrics() error
}
