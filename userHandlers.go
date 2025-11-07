package main

import (
	"net/http"
	"time"

	"github.com/gfrei/chirpy/internal/auth"
	"github.com/gfrei/chirpy/internal/database"
)

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

	resp := userJson{
		Id:          user.ID.String(),
		CreatedAt:   user.CreatedAt.GoString(),
		UpdatedAt:   user.UpdatedAt.GoString(),
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJson(w, http.StatusCreated, resp)
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, req *http.Request) {
	type jsonReq struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	params, err := decodeJson[jsonReq](req)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	hashedPassord, err := auth.HashPassword(params.Password)
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

	user, err := cfg.dbQueries.UpdateUser(req.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassord,
		UpdatedAt:      time.Now(),
		ID:             userId,
	})
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	resp := userJson{
		Id:          user.ID.String(),
		CreatedAt:   user.CreatedAt.GoString(),
		UpdatedAt:   user.UpdatedAt.GoString(),
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJson(w, http.StatusOK, resp)
}

func (cfg *apiConfig) refreshUserHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithJsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	refreshToken, err := cfg.dbQueries.GetRefreshToken(req.Context(), token)
	if err != nil || time.Now().Compare(refreshToken.ExpiresAt) > 0 || (refreshToken.RevokedAt.Valid && time.Now().Compare(refreshToken.RevokedAt.Time) > 0) {
		respondWithJsonError(w, http.StatusUnauthorized, "Not found")
		return
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secret)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	type jsonResponse struct {
		Token string `json:"token"`
	}

	resp := jsonResponse{
		Token: accessToken,
	}

	respondWithJson(w, http.StatusOK, resp)
}

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, req *http.Request) {
	type jsonParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	token, err := auth.MakeJWT(user.ID, cfg.secret)
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	_, err = cfg.dbQueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	type jsonResponse struct {
		Id           string `json:"id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		Email        string `json:"email"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	resp := jsonResponse{
		Id:           user.ID.String(),
		CreatedAt:    user.CreatedAt.GoString(),
		UpdatedAt:    user.UpdatedAt.GoString(),
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJson(w, http.StatusOK, resp)
}
