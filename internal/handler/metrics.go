package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rakett1337/pipedrive/internal/metrics"
)

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metricsCounts := metrics.GetMetricsCounts()

	metricsData := map[string]interface{}{}
	for key, count := range metricsCounts {
		method := key[:3]
		path := key[3:]
		metricsData[method+"_mean_duration"] = metrics.CalculateAverageDuration(method, path).String()
		metricsData[method+"_count"] = count
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metricsData)
}
