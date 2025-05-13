package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"user-service/api"
	"user-service/metrics"
	"user-service/routing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const defaultPort = 8080
const defaultMetricsPort = 8081
const defaultPprofPort = 8082

type App struct {
	mainPort     int
	metricsPort  int
	pprofPort    int
	pprofEnabled bool

	mainServer    *http.Server
	metricsServer *http.Server
	pprofServer   *http.Server
}

func New(opts ...AppOpt) *App {
	app := &App{
		mainPort:    defaultPort,
		metricsPort: defaultMetricsPort,
		pprofPort:   defaultPprofPort,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (s *App) Start() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	metrics.InitMetrics()

	s.mainServer = createMainServer(s.mainPort)
	s.metricsServer = createMetricsServer(s.metricsPort)

	startServer(s.mainServer, s.mainPort, "main")
	startServer(s.metricsServer, s.metricsPort, "metrics")

	if s.pprofEnabled {
		s.pprofServer = createPprofServer(s.pprofPort)
		startServer(s.pprofServer, s.pprofPort, "pprof")
	} else {
		slog.Info("pprof is disabled (ENABLE_PPROF != true)")
	}
}

func (s *App) Stop() {
	shutdown(s.mainServer, s.mainPort, "main")
	shutdown(s.metricsServer, s.metricsPort, "metrics")
	if s.pprofEnabled {
		shutdown(s.pprofServer, s.pprofPort, "pprof")
	}
}

func createMainServer(port int) *http.Server {
	mainMux := http.NewServeMux()
	routing.ConfigureMainRoutes(mainMux)

	mainServer := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: api.LoggingMiddleware(
			api.MetricsMiddleware(
				mainMux,
			),
		),
	}
	return mainServer
}

func createMetricsServer(port int) *http.Server {
	metricsMux := http.NewServeMux()
	metricsMux.Handle("GET /metrics", promhttp.Handler())

	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: metricsMux,
	}
	return metricsServer
}

func createPprofServer(port int) *http.Server {
	pprofMux := http.DefaultServeMux
	pprofServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: pprofMux,
	}
	return pprofServer
}

func startServer(s *http.Server, port int, name string) {
	go func() {
		slog.Info("server is starting...", "port", port, "name", name)
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error starting server", "error", err, "port", port)
			os.Exit(1)
		}
	}()
}

func shutdown(server *http.Server, port int, name string) {
	slog.Info("server is shutting down", "port", port, "name", name)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("error shutting down server", "error", err, "name", name)
		os.Exit(1)
	}

	slog.Info("server successfully shut down", "port", port, "name", name)
}
