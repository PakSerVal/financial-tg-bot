package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace     = "tgbot"
	httpSubSystem = "http"
	grpcSubSystem = "grpc"
)

var incomingMessageTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: httpSubSystem,
		Name:      "incoming_message_total",
	},
	[]string{"result"},
)

var messageProcessedTime = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: httpSubSystem,
		Name:      "message_processed_time",
		Buckets:   []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2},
	},
	[]string{"command"},
)

var serverProcessedTime = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: grpcSubSystem,
		Name:      "server_response_time_seconds",
		Buckets:   []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2},
	},
	[]string{"method"},
)

var clientProcessedTime = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: grpcSubSystem,
		Name:      "client_request_time_seconds",
		Buckets:   []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2},
	},
	[]string{"method"},
)

func IncomingMessageTotal(processedResult string) {
	incomingMessageTotal.WithLabelValues(processedResult).Inc()
}

func MessageProcessedTime(value float64, commandName string) {
	messageProcessedTime.WithLabelValues(commandName).Observe(value)
}

func GrpcServerResponseTime(value float64, methodName string) {
	serverProcessedTime.WithLabelValues(methodName).Observe(value)
}

func GrpcClientResponseTime(value float64, methodName string) {
	clientProcessedTime.WithLabelValues(methodName).Observe(value)
}
