package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

const remoteURL = "https://api.pipedrive.com/v1"

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/deals", dealsHandler)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", http.StatusNotFound)
	})

	handler := metricsMiddleware(r)
	log.Fatal(http.ListenAndServe(":80", handler))
}

func dealsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut:
		proxyHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL, err := url.Parse(remoteURL + r.URL.Path)

	// Copy query string from original request
	targetURL.RawQuery = r.URL.RawQuery
	client := &http.Client{}
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers from original request
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
