package monitoring

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"log/slog"
	"net/http"
)

var server *http.Server

var (
	NotificationsErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notifications_errors_total",
			Help: "Notifications Errors",
		},
		[]string{"message"},
	)
	NotificationsSuccessCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notifications_success_total",
			Help: "Notifications success",
		},
		[]string{"message"},
	)
)

func RegisterPrometheus() {
	prometheus.MustRegister(NotificationsErrorCount)
	prometheus.MustRegister(NotificationsSuccessCount)
}

func RunPrometheusServer(url string) {
	server = &http.Server{
		Addr: url,
	}

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start prometheus server: %v", err)
		}
	}()
}

func StopPrometheusServer(ctx context.Context) {
	if server != nil {
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Server forced to shutdown: %v", "error", err)
		}
		slog.Info("Server exited gracefully")
	} else {
		slog.Warn("Server is not running")
	}
}
