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
	connections   prometheus.GaugeFunc
	queueRecvLen  prometheus.GaugeFunc
	queueRecvCap  prometheus.GaugeFunc
	queueSendLen  prometheus.GaugeFunc
	queueSendCap  prometheus.GaugeFunc
	incomeEvents  *prometheus.CounterVec
	outcomeEvents *prometheus.CounterVec
	logger        logger.Logger
}

func newPromMetrics(logger logger.Logger, tr *wsTransport) *promMetrics {
	labels := []string{"event"}

	return &promMetrics{
		connections: metrics.CreateUsecaseGaugeFunc(usecase, "ws_connections", "Number of websocket connections", func() float64 {
			return float64(len(tr.connections))
		}),
		queueRecvLen: metrics.CreateUsecaseGaugeFunc(usecase, "ws_queue_recv_len", "Recv queue items", func() float64 {
			return float64(len(tr.busRecv))
		}),
		queueRecvCap: metrics.CreateUsecaseGaugeFunc(usecase, "ws_queue_recv_cap", "Recv queue length", func() float64 {
			return float64(cap(tr.busRecv))
		}),
		queueSendLen: metrics.CreateUsecaseGaugeFunc(usecase, "ws_queue_send_len", "Send queue items", func() float64 {
			return float64(len(tr.busSend))
		}),
		queueSendCap: metrics.CreateUsecaseGaugeFunc(usecase, "ws_queue_send_cap", "Send queue length", func() float64 {
			return float64(cap(tr.busSend))
		}),
		incomeEvents:  metrics.CreateUsecaseEventCounterVec(usecase, "ws_income_events", labels),
		outcomeEvents: metrics.CreateUsecaseEventCounterVec(usecase, "ws_outcome_events", labels),
		logger:        logger,
	}
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
