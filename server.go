package main

import (
	"net/http"
	"sync/atomic"

	"github.com/gfrei/chirpy/internal/database"
)

func newServer(dbQueries *database.Queries, platform, secret string) *http.Server {
	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
		platform:       platform,
		secret:         secret,
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

	mux.HandleFunc("GET /admin/metrics", apiCfg.fileserverHitsCountHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.fileserverHitsResetHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUserHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginUserHandler)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshUserHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeAccessTokenHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getAllChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.getChirpHandler)

	mux.HandleFunc("GET /api/healthz", readinessHandler)
}
