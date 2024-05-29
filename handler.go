package main

import (
	"encoding/json"
	"net/http"
)

func dealsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut:
		proxyHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func dealsMetricsHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	metrics := map[string]interface{}{
		"GET_mean_duration":  calculateAverageDuration(getMetrics).String(),
		"POST_mean_duration": calculateAverageDuration(postMetrics).String(),
		"PUT_mean_duration":  calculateAverageDuration(putMetrics).String(),
		"GET_count":          len(getMetrics),
		"POST_count":         len(postMetrics),
		"PUT_count":          len(putMetrics),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}
