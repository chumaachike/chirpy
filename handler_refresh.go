package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/chumaachike/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	// Extract refresh token from headers
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	// Look up user from refresh token (validates expiry + revocation in SQL)
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is invalid or expired", err)
		return
	}

	// Create new 1-hour access token (JWT)
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create new token", err)
		return
	}

	// Respond with new access token
	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Extract refresh token from headers
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	// Attempt to revoke
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil && err != sql.ErrNoRows {
		respondWithError(w, http.StatusInternalServerError, "Unable to revoke refresh token", err)
		return
	}

	// Success (idempotent): no content
	w.WriteHeader(http.StatusNoContent)
}
