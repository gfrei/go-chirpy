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
	w.Header().Set("Content-Type", "application/json")
	var dat []byte

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		dat, _ = getErrorJsonData("Something went wrong")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	nilParams := params == parameters{}

	if nilParams {
		dat, _ = getErrorJsonData("Something went wrong")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		dat, _ = getErrorJsonData("Chirp is too long")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	dat, _ = getOkJsonData()
	w.WriteHeader(200)
	w.Write(dat)
}

type errorJson struct {
	Error string `json:"error"`
}

func getErrorJsonData(message string) ([]byte, error) {
	respBody := errorJson{
		Error: message,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		return nil, err
	}

	return dat, nil
}

type okJson struct {
	Valid bool `json:"valid"`
}

func getOkJsonData() ([]byte, error) {
	respBody := okJson{
		Valid: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		return nil, err
	}

	return dat, nil
}
