package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "tgbot"
	subSystem = "http"
)

var incomingMessageTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subSystem,
		Name:      "incoming_message_total",
	},
	[]string{"result"},
)

var messageProcessedTime = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subSystem,
		Name:      "message_processed_time",
		Buckets:   []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2},
	},
	[]string{"command"},
)

func IncomingMessageTotal(processedResult string) {
	incomingMessageTotal.WithLabelValues(processedResult).Inc()
}

func MessageProcessedTime(value float64, commandName string) {
	messageProcessedTime.WithLabelValues(commandName).Observe(value)
}
