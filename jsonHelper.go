package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func decodeJson[T any](req *http.Request) (T, error) {
	var params T
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&params); err != nil {
		var zero T
		return zero, fmt.Errorf("invalid JSON: %w", err)
	}

	return params, nil
}
