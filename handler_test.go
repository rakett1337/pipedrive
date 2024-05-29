package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var apiToken = os.Getenv("API_TOKEN")

func TestDealsHandler(t *testing.T) {
	testCases := []struct {
		name           string
		authMethod     string
		expectedStatus int
	}{
		{"Header Authentication", "header", http.StatusOK},
		{"Query Parameter Authentication", "query", http.StatusOK},
		{"Unauthorized", "", http.StatusUnauthorized},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/deals", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tc.authMethod == "header" {
				req.Header.Set("X-API-Token", apiToken)
			} else if tc.authMethod == "query" {
				req.URL.RawQuery = "api_token=" + apiToken
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(dealsHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}
		})
	}
}

func TestDealsMetricsHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/deals", dealsHandler)
	mux.HandleFunc("/metrics", dealsMetricsHandler)
	loggedMux := metricsMiddleware(mux)
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

	metricsSlices := []struct {
		name  string
		slice []metrics
	}{
		{"GET", getMetrics},
		{"POST", postMetrics},
		{"PUT", putMetrics},
	}

	for _, m := range metricsSlices {
		t.Run(m.name+" Count", func(t *testing.T) {
			if len(m.slice) != 1 {
				t.Errorf("unexpected count for %s metrics: got %d, want 1", m.name, len(m.slice))
			}
		})
	}
}
