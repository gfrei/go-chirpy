package main

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gfrei/chirpy/internal/auth"
	"github.com/gfrei/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, req *http.Request) {
	chirpId := req.PathValue("id")

	chirpUUID, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), chirpUUID)
	if err != nil {
		respondWithJsonError(w, http.StatusNotFound, "Not found")
		return
	}

	chirpJson := getChirpJson(chirp)

	respondWithJson(w, http.StatusOK, chirpJson)
}

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, req *http.Request) {
	chirpId := req.PathValue("chirpID")

	chirpUUID, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), chirpUUID)
	if err != nil {
		respondWithJsonError(w, http.StatusNotFound, "Not found")
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithJsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithJsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if chirp.UserID != userId {
		respondWithJsonError(w, http.StatusForbidden, "Forbidden")
		return
	}

	err = cfg.dbQueries.DeleteChirp(req.Context(), chirpUUID)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, req *http.Request) {
	authorId := req.URL.Query().Get("author_id")

	var chirps []database.Chirp
	var err error

	if authorId != "" {
		authorUUID, err := uuid.Parse(authorId)
		if err != nil {
			respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
			return
		}

		chirps, err = cfg.dbQueries.GetAllChirpsFromUser(req.Context(), authorUUID)
		if err != nil {
			respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
			return
		}
	} else {
		chirps, err = cfg.dbQueries.GetAllChirps(req.Context())
		if err != nil {
			respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
			return
		}
	}

	sorting := req.URL.Query().Get("sort")
	if sorting == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Compare(chirps[j].CreatedAt) > 0 })
	}

	chirpsJson := make([]chirpJson, 0)

	for _, chirp := range chirps {
		chirpJson := getChirpJson(chirp)
		chirpsJson = append(chirpsJson, chirpJson)
	}

	respondWithJson(w, http.StatusOK, chirpsJson)
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, req *http.Request) {
	type jsonReq struct {
		Body string `json:"body"`
	}

	params, err := decodeJson[jsonReq](req)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithJsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithJsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if len(params.Body) == 0 {
		respondWithJsonError(w, http.StatusBadRequest, fmt.Sprintf("Something went wrong: %v", err))
		return
	}

	if len(params.Body) > 140 {
		respondWithJsonError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		UserID: userId,
		Body:   params.Body,
	})
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, fmt.Sprintf("Something went wrong: %v", err))
		return
	}

	chirpJson := getChirpJson(chirp)

	respondWithJson(w, http.StatusCreated, chirpJson)
}
