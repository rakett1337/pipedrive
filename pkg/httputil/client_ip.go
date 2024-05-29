package httputil

import "net/http"

// GetClientIP returns the client IP address from the request.
func GetClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
