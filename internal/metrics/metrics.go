package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	Status    int
	Duration  time.Duration
	Method    string
	Path      string
	IPAddress string
	UserAgent string
}

var (
	metricsData sync.Map
	mutex       sync.Mutex
)

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

func GetMetrics(method, path string) []Metrics {
	key := method + path
	if v, ok := metricsData.Load(key); ok {
		return v.([]Metrics)
	}
	return []Metrics{}
}

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

func GetMetricsCounts() map[string]int {
	counts := make(map[string]int)
	metricsData.Range(func(key, value interface{}) bool {
		metrics := value.([]Metrics)
		counts[key.(string)] = len(metrics)
		return true
	})
	return counts
}
