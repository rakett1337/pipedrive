package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type Metrics struct {
	Duration time.Duration
	Method   string
	Path     string
}

var (
	getMetrics  []Metrics
	postMetrics []Metrics
	putMetrics  []Metrics
	mutex       sync.Mutex
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		if r.URL.Path == "/deals" {
			metrics := Metrics{
				Duration: duration,
				Method:   r.Method,
				Path:     r.URL.Path,
			}
			saveMetrics(metrics)
		}
		log.Printf("%s %s %s %s", r.Method, r.URL.Path, duration, r.UserAgent())
	})
}

func saveMetrics(metrics Metrics) {
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
