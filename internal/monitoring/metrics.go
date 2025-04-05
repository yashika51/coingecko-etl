package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	FetchSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "coingecko_fetch_success_total",
			Help: "Total number of successful API fetches",
		})

	FetchFailure = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "coingecko_fetch_failure_total",
			Help: "Total number of failed API fetches",
		})

	RecordsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "coingecko_records_processed_total",
			Help: "Total number of coin records processed",
		})
)

func Init() {
	prometheus.MustRegister(FetchSuccess, FetchFailure, RecordsProcessed)
}
