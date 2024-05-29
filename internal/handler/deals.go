package handler

import (
	"io"
	"net/http"
	"net/url"
)

const remoteURL = "https://api.pipedrive.com/v1"

// DealsHandler is an HTTP handler that proxies requests to the Pipedrive API.
// It supports GET, POST, and PUT methods. Returns StatusMethodNotAllowed for other methods.
// The handler is used to retrieve, create, and update deals.
func DealsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut:
		proxyHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL, err := url.Parse(remoteURL + r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	targetURL.RawQuery = r.URL.RawQuery

	client := &http.Client{}
	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	proxyReq.Header = r.Header

	proxyResp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Failed to perform request", http.StatusInternalServerError)
		return
	}
	defer proxyResp.Body.Close()

	for key, values := range proxyResp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(proxyResp.StatusCode)
	if _, err := io.Copy(w, proxyResp.Body); err != nil {
		http.Error(w, "Failed to copy response body", http.StatusInternalServerError)
	}
}
