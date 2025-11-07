package main

import (
	"net/http"
	"time"

	"github.com/gfrei/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhookSetUserRedHandler(w http.ResponseWriter, req *http.Request) {
	type jsonReq struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	params, err := decodeJson[jsonReq](req)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	uuid, err := uuid.Parse(params.Data.UserId)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	user, err := cfg.dbQueries.GetUserById(req.Context(), uuid)
	if err != nil {
		respondWithJsonError(w, http.StatusNotFound, "User not found")
		return
	}

	if user.IsChirpyRed {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.dbQueries.UpdateUserIsChirpyRed(req.Context(), database.UpdateUserIsChirpyRedParams{
		IsChirpyRed: true,
		ID:          uuid,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
