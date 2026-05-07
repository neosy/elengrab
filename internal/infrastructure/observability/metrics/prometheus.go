package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	dbTableRows = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sqlite_table_rows",
			Help: "Number of rows in SQLite tables",
		},
		[]string{"db_name", "table"},
	)
)

func Register() {
	prometheus.MustRegister(
		httpRequestsTotal, httpRequestDuration,
		dbTableRows,
	)
}
