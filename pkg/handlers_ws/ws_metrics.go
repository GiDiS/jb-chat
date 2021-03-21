package handlers_ws

import (
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/GiDiS/jb-chat/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const usecase = "WebsocketTransport"

type Metrics interface {
	SetConnections(cnt int)
	IncIncome(event string)
	IncOutcome(event string)
}

type promMetrics struct {
	Metrics
	connections   *prometheus.Gauge
	incomeEvents  *prometheus.CounterVec
	outcomeEvents *prometheus.CounterVec
	logger        logger.Logger
}

func newPromMetrics(logger logger.Logger) *promMetrics {
	labels := []string{"event"}

	connections := metrics.CreateUsecaseGauge(usecase, "connections", "Number of websocket connections")

	metrics := promMetrics{
		connections:   &connections,
		incomeEvents:  metrics.CreateUsecaseEventCounterVec(usecase, "income_events", labels),
		outcomeEvents: metrics.CreateUsecaseEventCounterVec(usecase, "outcome_events", labels),
		logger:        logger,
	}

	return &metrics
}

func (m promMetrics) SetConnections(cnt int) {
	(*m.connections).Set(float64(cnt))
}

func (m promMetrics) IncIncome(event string) {
	metric, err := m.incomeEvents.GetMetricWithLabelValues(event)
	if err != nil {
		m.logger.WithError(err).WithField("action", event).Errorf("metric income event not found")
	} else {
		metric.Inc()
	}
}

func (m promMetrics) IncOutcome(event string) {
	metric, err := m.outcomeEvents.GetMetricWithLabelValues(event)
	if err != nil {
		m.logger.WithError(err).WithField("action", event).Errorf("metric outcome event not found")
	} else {
		metric.Inc()
	}
}
