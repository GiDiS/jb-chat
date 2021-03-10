package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

func TrackDuration(histogram prometheus.Histogram, start time.Time) {
	histogram.Observe(time.Since(start).Seconds())
}
