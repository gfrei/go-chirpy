package main

import (
	"net/http"
	"sync/atomic"
)

func newServer() *http.Server {
	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	addHandlers(mux, apiCfg)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	return server
}

func addHandlers(mux *http.ServeMux, apiCfg *apiConfig) {
	fsHandlerWithMetrics := apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.Handle("/app/", fsHandlerWithMetrics)
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /api/metrics", apiCfg.fileserverHitsCountHandler)
	mux.HandleFunc("POST /api/reset", apiCfg.fileserverHitsResetHandler)
}
