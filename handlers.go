package main

import (
	"encoding/json"
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
	type parameters struct {
		Body string `json:"body"`
	}

	params := parameters{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&params); err != nil {
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

func respondWithJsonError(w http.ResponseWriter, code int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(w, code, errorResponse{Error: message})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	dat, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong"})
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}
