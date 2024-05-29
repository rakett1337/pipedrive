package main

import (
	"log"
	"net/http"

	"github.com/rakett1337/pipedrive/internal/handler"
	"github.com/rakett1337/pipedrive/internal/middleware"
)

func main() {
	log.SetFlags(log.LstdFlags)
	mux := http.NewServeMux()
	mux.HandleFunc("/deals", handler.DealsHandler)
	mux.HandleFunc("/metrics", handler.MetricsHandler)

	loggedMux := middleware.MetricsMiddleware(mux)

	log.Println("Starting server on :80")
	log.Fatal(http.ListenAndServe(":80", loggedMux))
}
