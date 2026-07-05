package middlewares

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var HttpRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Nombre total de requêtes HTTP traitées",
	},
	[]string{"status", "version"},
)

func PrometheusMiddleware(version string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{w, http.StatusOK}
		next(lrw, r)

		statusCode := lrw.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		statusLabel := "success"
		if statusCode >= 500 {
			statusLabel = "error"
		}

		HttpRequestsTotal.WithLabelValues(statusLabel, version).Inc()
	}
}

// utilitaire pour intercepter le Status Code HTTP

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}