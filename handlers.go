package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gfrei/chirpy/internal/auth"
	"github.com/gfrei/chirpy/internal/database"
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

func (cfg *apiConfig) revokeAccessTokenHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithJsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	refreshToken, err := cfg.dbQueries.GetRefreshToken(req.Context(), token)
	if err != nil || time.Now().Compare(refreshToken.ExpiresAt) > 0 {
		respondWithJsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(req.Context(), database.RevokeRefreshTokenParams{
		Token:     refreshToken.Token,
		UpdatedAt: time.Now(),
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		respondWithJsonError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
