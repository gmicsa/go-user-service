package api

import (
	"net/http"
	"regexp"
	"time"
	"user-service/metrics"

	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel/metric"
)

// Regex pattern to replace numbers in the path
var numberPattern = regexp.MustCompile(`\d+`)

func normalizePath(path string) string {
	return numberPattern.ReplaceAllString(path, "{id}")
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := normalizePath(r.URL.Path)
		attributes := metric.WithAttributes(attribute.String("url", path))
		metrics.RequestCounter.Add(r.Context(), 1, attributes)

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()

		metrics.RequestLatency.Record(r.Context(), duration, attributes)
	})
}
