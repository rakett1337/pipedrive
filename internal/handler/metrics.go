package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rakett1337/pipedrive/internal/metrics"
)

// MetricsHandler is an HTTP handler that retrieves and returns the metrics data collected by MetricsMiddleware.
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metricsCounts := metrics.GetMetricsCounts()
	metricsData := map[string]interface{}{}
	for key, count := range metricsCounts {
		parts := strings.SplitN(key, "/", 2)
		if len(parts) == 2 {
			method := parts[0]
			path := "/" + parts[1]
			metricsData[method+"_path"] = path
			metricsData[method+"_avg_duration"] = metrics.CalculateAverageDuration(method, path).String()
			metricsData[method+"_count"] = count
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metricsData)
}
