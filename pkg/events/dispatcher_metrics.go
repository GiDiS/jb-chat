package events

import (
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/GiDiS/jb-chat/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const usecase = "EventsDispatcher"

type Metrics interface {
	RegisterEvent(event string)
	TrackEvent(event string) Stopper
	IncError(event string)
	IncUnhandled(event string)
	IncHandled(event string)
	IncNotify(event string)
}

type promMetrics struct {
	Metrics
	eventsRuns    *prometheus.CounterVec
	eventsErrors  *prometheus.CounterVec
	eventsTimes   *prometheus.HistogramVec
	unhandledRuns *prometheus.CounterVec
	handledRuns   *prometheus.CounterVec
	notifyRuns    *prometheus.CounterVec
	logger        logger.Logger
}

func newPromMetrics(logger logger.Logger) *promMetrics {
	labels := []string{"event"}

	metrics := promMetrics{
		eventsRuns:    metrics.CreateUsecaseEventCounterVec(usecase, "exec_event_total", labels),
		eventsTimes:   metrics.CreateUsecaseEventTimingHistogramVec(usecase, labels),
		eventsErrors:  metrics.CreateUsecaseEventErrorCounterVec(usecase, labels),
		unhandledRuns: metrics.CreateUsecaseEventCounterVec(usecase, "unhandled_events", labels),
		handledRuns:   metrics.CreateUsecaseEventCounterVec(usecase, "handled_events", labels),
		notifyRuns:    metrics.CreateUsecaseEventCounterVec(usecase, "notify_events", labels),
		logger:        logger,
	}

	return &metrics
}

type Stopper func()

func (m promMetrics) RegisterEvent(event string) {
	if _, err := m.eventsRuns.GetMetricWithLabelValues(event); err != nil {
		m.logger.WithError(err).Errorf("Metric registering failed: %s", event)
	}
	if _, err := m.eventsTimes.GetMetricWithLabelValues(event); err != nil {
		m.logger.WithError(err).Errorf("Metric registering failed: %s", event)
	}
	if _, err := m.eventsErrors.GetMetricWithLabelValues(event); err != nil {
		m.logger.WithError(err).Errorf("Metric registering failed: %s", event)
	}
}

func (m promMetrics) TrackEvent(event string) Stopper {
	start := time.Now()
	metric, err := m.eventsRuns.GetMetricWithLabelValues(event)
	if err != nil {
		m.logger.WithError(err).WithField("event", event).Errorf("metric event runs not found")
	} else if metric != nil {
		metric.Inc()
	}

	return func() {
		metric, err := m.eventsTimes.GetMetricWithLabelValues(event)
		if err != nil {
			m.logger.WithError(err).WithField("event", event).Errorf("metric event times not found")
		} else if metric != nil {
			metric.Observe(time.Since(start).Seconds())
		}
	}
}

func (m promMetrics) IncError(event string) {
	metric, err := m.eventsErrors.GetMetricWithLabelValues(event)
	if err != nil {
		m.logger.WithError(err).WithField("action", event).Errorf("metric event errors not found")
	} else {
		metric.Inc()
	}
}

func (m promMetrics) IncUnhandled(event string) {
	metric, err := m.unhandledRuns.GetMetricWithLabelValues(event)
	if err != nil {
		m.logger.WithError(err).WithField("action", event).Errorf("metric event unhandled not found")
	} else {
		metric.Inc()
	}
}

func (m promMetrics) IncHandled(event string) {
	metric, err := m.handledRuns.GetMetricWithLabelValues(event)
	if err != nil {
		m.logger.WithError(err).WithField("action", event).Errorf("metric event unhandled not found")
	} else {
		metric.Inc()
	}
}

func (m promMetrics) IncHandledNotify(event string) {
	metric, err := m.notifyRuns.GetMetricWithLabelValues(event)
	if err != nil {
		m.logger.WithError(err).WithField("action", event).Errorf("metric event notify not found")
	} else {
		metric.Inc()
	}
}
