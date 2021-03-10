package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const Ns = "jb-chat"

func CreateCounter(name, help string, labels map[string]string) prometheus.Counter {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   Ns,
		Name:        name,
		ConstLabels: prometheus.Labels(labels),
		Help:        help,
	})
	prometheus.MustRegister(counter)
	return counter
}

func CreateUsecaseGauge(usecaseName, gaugeName, gaugeHelp string) prometheus.Gauge {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Ns,
		Subsystem: "usecase",
		Name:      gaugeName,
		ConstLabels: prometheus.Labels{
			"usecase": usecaseName,
		},
		Help: gaugeHelp,
	})
	return gauge
}

func CreateUsecaseCounter(usecaseName string) prometheus.Counter {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Ns,
		Subsystem: "usecase",
		Name:      "exec_total",
		ConstLabels: prometheus.Labels{
			"usecase": usecaseName,
		},
		Help: "Number of executions of usecase",
	})

	prometheus.MustRegister(counter)
	return counter
}

func CreateUsecaseActionCounterVec(usecaseName string, labelNames []string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Ns,
		Subsystem: "usecase",
		Name:      "exec_action_total",
		ConstLabels: prometheus.Labels{
			"usecase": usecaseName,
		},
		Help: "Number of action executions of usecase",
	}, labelNames)

	prometheus.MustRegister(*counter)
	return counter
}

func CreateUsecaseActionErrorsCounterVec(usecaseName string, labelNames []string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Ns,
		Subsystem: "usecase",
		Name:      "errors_action_total",
		ConstLabels: prometheus.Labels{
			"usecase": usecaseName,
		},
		Help: "Number of action fails of usecase",
	}, labelNames)

	prometheus.MustRegister(*counter)
	return counter
}

func CreateUsecaseErrorCounter(usecaseName, errorCode string) prometheus.Counter {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Ns,
		Subsystem: "usecase",
		Name:      "errors_total",
		Help:      "Number of errors in usecase",
		ConstLabels: prometheus.Labels{
			"usecase": usecaseName,
			"error":   errorCode,
		},
	})
	prometheus.MustRegister(counter)
	return counter
}

func CreateUsecaseTimingHistogram(usecaseName string) prometheus.Histogram {
	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: Ns,
		Subsystem: "usecase",
		Name:      "exec_seconds",
		ConstLabels: prometheus.Labels{
			"usecase": usecaseName,
		},
		Help: "Usecase run duration in seconds",
	})
	prometheus.MustRegister(histogram)
	return histogram
}

func CreateUsecaseActionTimingHistogramVec(usecaseName string, labelNames []string) *prometheus.HistogramVec {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: Ns,
		Subsystem: "usecase",
		Name:      "exec_action_seconds",
		ConstLabels: prometheus.Labels{
			"usecase": usecaseName,
		},
		Help: "Usecase action run duration in seconds",
	}, labelNames)
	prometheus.MustRegister(*histogram)
	return histogram
}
