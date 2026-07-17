package core_metrics

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	core_http_middleware "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/middleware"
	core_http_response "github.com/rallaverdi/golang-todoapp/internal/core/transport/http/response"
)

type Metrics struct {
	registry         *prometheus.Registry
	requestsTotal    *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	requestsInFlight prometheus.Gauge
}

func NewMetrics() *Metrics {
	registry := prometheus.NewRegistry()

	metrics := &Metrics{
		registry: registry,
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "todoapp",
				Subsystem: "http",
				Name:      "requests_total",
				Help:      "Total number of processed HTTP requests.",
			},
			[]string{"version", "method", "route", "status"},
		),

		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "todoapp",
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "HTTP request duration in seconds.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"version", "method", "route", "status"},
		),

		requestsInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "todoapp",
				Subsystem: "http",
				Name:      "requests_in_flight",
				Help:      "Number of requests currently in flight.",
			}),
	}

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		metrics.requestsTotal,
		metrics.requestDuration,
		metrics.requestsInFlight,
	)

	return metrics
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

func (m *Metrics) HTTPMiddleware(apiVersion string) core_http_middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()
			rw := core_http_response.NewResponseWriter(w)

			m.requestsInFlight.Inc()
			defer m.requestsInFlight.Dec()

			defer func() {
				statusCode := rw.GetStatusCode()

				if recovered := recover(); recovered != nil {
					statusCode = http.StatusInternalServerError
					m.observeRequest(apiVersion, r, statusCode, startedAt)
					panic(recovered)
				}

				m.observeRequest(apiVersion, r, statusCode, startedAt)
			}()

			next.ServeHTTP(rw, r)
		})
	}
}

func (m *Metrics) observeRequest(
	apiVersion string,
	r *http.Request,
	statusCode int,
	startedAt time.Time,
) {
	route := strings.TrimPrefix(r.Pattern, r.Method+" ")
	if route == "" {
		route = "unknown"
	}

	labels := []string{
		apiVersion,
		r.Method,
		route,
		strconv.Itoa(statusCode),
	}

	m.requestsTotal.WithLabelValues(labels...).Inc()
	m.requestDuration.WithLabelValues(labels...).Observe(time.Since(startedAt).Seconds())
}
