package metrics

import (
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	otelSDK "go.opentelemetry.io/otel/sdk/metric"
)

var (
	RequestCounter metric.Int64Counter
	RequestLatency metric.Float64Histogram
)

func InitMetrics() {
	exporter, err := prometheus.New()
	if err != nil {
		slog.Error("Failed to create Prometheus exporter", "error", err)
		os.Exit(1)
	}

	provider := otelSDK.NewMeterProvider(
		otelSDK.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	meter := provider.Meter("user_service")

	RequestCounter, err = meter.Int64Counter(
		"request_total",
		metric.WithDescription("Total number of incoming requests"),
	)
	if err != nil {
		slog.Error("Failed to create requestCounter", "error", err)
		os.Exit(1)
	}

	RequestLatency, err = meter.Float64Histogram(
		"request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests in seconds"),
	)
	if err != nil {
		slog.Error("Failed to create requestLatency", "error", err)
		os.Exit(1)
	}
}
