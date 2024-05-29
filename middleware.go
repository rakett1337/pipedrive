package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type metrics struct {
	Status    int
	Duration  time.Duration
	Method    string
	IPAddress string
	UserAgent string
	Headers   http.Header
	Query     map[string][]string
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}

var (
	getMetrics  []metrics
	postMetrics []metrics
	putMetrics  []metrics
	mutex       sync.Mutex
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		ip := getClientIP(r)

		if r.URL.Path == "/deals" {
			metrics := metrics{
				Status:    rw.status,
				Duration:  duration,
				Method:    r.Method,
				IPAddress: ip,
				UserAgent: r.UserAgent(),
				Headers:   r.Header,
				Query:     r.URL.Query(),
			}
			savemetrics(metrics)
		}
		testMode := os.Getenv("TEST_MODE") == "true"
		if !testMode {
			log.Printf("%d %s %s %s %s %s", rw.status, r.Method, r.URL.Path, duration, r.UserAgent(), ip)
		}
	})
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func savemetrics(metrics metrics) {
	mutex.Lock()
	defer mutex.Unlock()

	switch metrics.Method {
	case http.MethodGet:
		getMetrics = append(getMetrics, metrics)
	case http.MethodPost:
		postMetrics = append(postMetrics, metrics)
	case http.MethodPut:
		putMetrics = append(putMetrics, metrics)
	}
}

func calculateAverageDuration(metrics []metrics) time.Duration {
	totalDuration := time.Duration(0)
	for _, metric := range metrics {
		totalDuration += metric.Duration
	}
	if len(metrics) == 0 {
		return 0
	}
	return totalDuration / time.Duration(len(metrics))
}
