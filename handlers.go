package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithJsonError(w, 400, "Something went wrong")
		return
	}

	nilParams := params == parameters{}

	if nilParams {
		respondWithJsonError(w, 400, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithJsonError(w, 400, "Chirp is too long")
		return
	}

	type jsonResponse struct {
		Valid bool `json:"valid"`
	}

	respondWithJson(w, http.StatusOK, jsonResponse{
		Valid: true,
	})
}

func respondWithJsonError(w http.ResponseWriter, responseValue int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	myError := errorResponse{
		Error: message,
	}

	respondWithJson(w, responseValue, myError)
}

func respondWithJson(w http.ResponseWriter, responseValue int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	code := responseValue

	dat, err := json.Marshal(payload)
	if err != nil {
		code = 500
	}

	w.WriteHeader(code)
	w.Write(dat)

	return err
}
