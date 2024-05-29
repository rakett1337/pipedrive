package handler

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rakett1337/pipedrive/internal/metrics"
	"github.com/rakett1337/pipedrive/internal/middleware"
)

func TestDealsMetricsHandler(t *testing.T) {
	log.SetOutput(io.Discard)
	mux := http.NewServeMux()
	mux.HandleFunc("/deals", DealsHandler)
	mux.HandleFunc("/metrics", MetricsHandler)
	loggedMux := middleware.MetricsMiddleware(mux)
	testServer := httptest.NewServer(loggedMux)
	defer testServer.Close()

	requestMethods := []string{"GET", "POST", "PUT"}
	for _, method := range requestMethods {
		req, err := http.NewRequest(method, testServer.URL+"/deals", nil)
		if err != nil {
			t.Fatal(err)
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
	}

	metricsCounts := metrics.GetMetricsCounts()
	for _, method := range requestMethods {
		t.Run(method+" Count", func(t *testing.T) {
			key := method + "/deals"
			if count, ok := metricsCounts[key]; !ok || count != 1 {
				t.Errorf("unexpected count for %s metrics: got %d, want 1", method, count)
			}
		})
	}
}
