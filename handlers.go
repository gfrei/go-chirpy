package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gfrei/chirpy/internal/auth"
	"github.com/gfrei/chirpy/internal/database"
	"github.com/google/uuid"
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

	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.DeleteAllUsers(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Something went wrong"))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintln("Reset Server")))
}

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

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(req.Context())
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
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
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	params, err := decodeJson[jsonReq](req)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	_, err = auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithJsonError(w, http.StatusUnauthorized, "Unauthorized")
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

	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		UserID: params.UserId,
		Body:   params.Body,
	})
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	chirpJson := getChirpJson(chirp)

	respondWithJson(w, http.StatusCreated, chirpJson)
}

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, req *http.Request) {
	type jsonParams struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	params, err := decodeJson[jsonParams](req)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithJsonError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithJsonError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	expiresInSeconds := params.ExpiresInSeconds
	if expiresInSeconds == 0 || expiresInSeconds > 3600 {
		expiresInSeconds = 3600
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Second*time.Duration(expiresInSeconds))
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	type jsonResponse struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
		Token     string `json:"token"`
	}

	resp := jsonResponse{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt.GoString(),
		UpdatedAt: user.UpdatedAt.GoString(),
		Email:     user.Email,
		Token:     token,
	}

	respondWithJson(w, http.StatusOK, resp)
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request) {
	type jsonParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params, err := decodeJson[jsonParams](req)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	hashedPassord, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassord,
	})
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	type jsonResponse struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}

	resp := jsonResponse{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt.GoString(),
		UpdatedAt: user.UpdatedAt.GoString(),
		Email:     user.Email,
	}

	respondWithJson(w, http.StatusCreated, resp)
}
