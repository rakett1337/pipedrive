package main

import (
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags)
	mux := http.NewServeMux()
	mux.HandleFunc("/deals", dealsHandler)
	mux.HandleFunc("/metrics", dealsMetricsHandler)

	loggedMux := metricsMiddleware(mux)

	log.Println("Starting server on :80")
	log.Fatal(http.ListenAndServe(":80", loggedMux))
}
