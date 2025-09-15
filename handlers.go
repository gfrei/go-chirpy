package main

import (
	"fmt"
	"net/http"

	"github.com/gfrei/chirpy/internal/stringvalidator"
)

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) fileserverHitsCountHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) fileserverHitsResetHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	w.Write([]byte(fmt.Sprintf("Reset Hits: %v", cfg.fileserverHits.Load())))
}

func validateChirpHandler(w http.ResponseWriter, req *http.Request) {
	params, err := decodeJson[struct {
		Body string `json:"body"`
	}](req)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(params.Body) == 0 {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithJsonError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	type jsonResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	respondWithJson(w, http.StatusOK, jsonResponse{CleanedBody: stringvalidator.StatelessClean(params.Body, []string{"kerfuffle", "sharbert", "fornax"})})
}
