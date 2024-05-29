package httputil

import "net/http"

// ResponseWriter is a wrapper around http.ResponseWriter that captures the HTTP status code.
type ResponseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader writes the HTTP status code to the response and saves it in the ResponseWriter.
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write writes the response body to the ResponseWriter.
// If no status code has been written, it defaults to http.StatusInternalServerError (500).
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusInternalServerError
	}
	return rw.ResponseWriter.Write(b)
}

// Status returns the HTTP status code written to the ResponseWriter.
func (rw *ResponseWriter) Status() int {
	return rw.status
}
