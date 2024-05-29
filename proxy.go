package main

import (
	"io"
	"net/http"
	"net/url"
)

const remoteURL = "https://api.pipedrive.com/v1"

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
