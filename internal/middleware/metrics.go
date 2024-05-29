package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/rakett1337/pipedrive/internal/metrics"
	"github.com/rakett1337/pipedrive/pkg/httputil"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &httputil.ResponseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)
		duration := time.Since(start)

		metrics.SaveMetrics(metrics.Metrics{
			Status:    rw.Status(),
			Duration:  duration,
			Method:    r.Method,
			Path:      r.URL.Path,
			IPAddress: httputil.GetClientIP(r),
			UserAgent: r.UserAgent(),
		})

		log.Printf("%d %s %s %s %s %s", rw.Status(), r.Method, r.URL.Path, duration, r.UserAgent(), httputil.GetClientIP(r))
	})
}
