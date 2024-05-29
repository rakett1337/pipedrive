// Package metrics provides functionality for capturing and retrieving metrics data.
package metrics

import (
	"sync"
	"time"
)

// Metrics represents the structure of the captured metrics data.
type Metrics struct {
	Status    int           // HTTP status code
	Duration  time.Duration // Request duration
	Method    string        // HTTP method
	Path      string        // Request path
	IPAddress string        // Client IP address
	UserAgent string        // Client User agent
}

var (
	metricsData sync.Map
	mutex       sync.Mutex
)

// SaveMetrics saves the provided Metrics instance to the metrics data store.
func SaveMetrics(m Metrics) {
	mutex.Lock()
	defer mutex.Unlock()

	key := m.Method + m.Path
	var metrics []Metrics
	if v, ok := metricsData.Load(key); ok {
		metrics = v.([]Metrics)
	}
	metrics = append(metrics, m)
	metricsData.Store(key, metrics)
}

// GetMetrics retrieves the metrics data for the specified method and path.
func GetMetrics(method, path string) []Metrics {
	key := method + path
	if v, ok := metricsData.Load(key); ok {
		return v.([]Metrics)
	}
	return []Metrics{}
}

// CalculateAverageDuration calculates the average duration for the specified method and path.
func CalculateAverageDuration(method, path string) time.Duration {
	metrics := GetMetrics(method, path)
	totalDuration := time.Duration(0)
	for _, metric := range metrics {
		totalDuration += metric.Duration
	}
	if len(metrics) == 0 {
		return 0
	}
	return totalDuration / time.Duration(len(metrics))
}

// GetMetricsCounts retrieves the count of metrics for each unique method and path combination.
func GetMetricsCounts() map[string]int {
	counts := make(map[string]int)
	metricsData.Range(func(key, value interface{}) bool {
		metrics := value.([]Metrics)
		counts[key.(string)] = len(metrics)
		return true
	})
	return counts
}
