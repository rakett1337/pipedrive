package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var apiToken = os.Getenv("API_TOKEN")

func TestDealsHandler(t *testing.T) {
	if apiToken == "" {
		t.Error("API_TOKEN is not set")
		return
	}

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
			handler := http.HandlerFunc(DealsHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}
		})
	}
}
