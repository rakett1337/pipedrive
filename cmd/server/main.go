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
	mux.HandleFunc("GET /deals", handler.DealsHandler)
	mux.HandleFunc("POST /deals", handler.DealsHandler)
	mux.HandleFunc("PUT /deals/{id}", handler.DealsHandler)
	mux.HandleFunc("GET /metrics", handler.MetricsHandler)

	loggedMux := middleware.MetricsMiddleware(mux)

	log.Println("Starting server on :80")
	log.Fatal(http.ListenAndServe(":80", loggedMux))
}
