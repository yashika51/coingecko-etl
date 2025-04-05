package main

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	go http.ListenAndServe(":8080", nil)
}
