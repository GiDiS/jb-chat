package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type BaseMetrics interface {
	Run()
	Stop()
	DbError(error)
}

type BaseErpMetrics interface {
	ErpError(error)
}

type BasePrometheusMetrics struct {
	start        time.Time
	timing       prometheus.Histogram
	runCounter   prometheus.Counter
	errDBCounter prometheus.Counter
}

type BaseErpPrometheusMetrics struct {
	errERPCounter prometheus.Counter
}

func CreateBasePrometheusMetrics(usecase string) BasePrometheusMetrics {
	return BasePrometheusMetrics{
		runCounter:   CreateUsecaseCounter(usecase),
		timing:       CreateUsecaseTimingHistogram(usecase),
		errDBCounter: CreateUsecaseErrorCounter(usecase, "db"),
	}
}

func (p BasePrometheusMetrics) Run() {
	p.start = time.Now()
	p.runCounter.Inc()
}

func (p BasePrometheusMetrics) Stop() {
	TrackDuration(p.timing, p.start)
}

func (p BasePrometheusMetrics) DbError(err error) {
	p.errDBCounter.Inc()
}

func CreateBaseErpPrometheusMetrics(usecase string) BaseErpPrometheusMetrics {
	return BaseErpPrometheusMetrics{
		errERPCounter: CreateUsecaseErrorCounter(usecase, "erp"),
	}
}

func (b BaseErpPrometheusMetrics) ErpError(error) {
	b.errERPCounter.Inc()
}
