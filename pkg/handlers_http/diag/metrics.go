package diag

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func registerMetrics(r *mux.Router) {
	r.Handle("/metrics", promhttp.Handler())
}
