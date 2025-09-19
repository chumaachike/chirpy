package main

import (
	"database/sql"
	"net/http"

	"github.com/chumaachike/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDelete(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from JWT
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bad token format", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	// Parse chirpID from path
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id", err)
		return
	}

	// Fetch chirp from DB
	chirp, err := cfg.db.GetChrirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "chirp not found", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "could not fetch chirp", err)
		}
		return
	}

	// Check ownership
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "you cannot delete someone else's chirp", nil)
		return
	}

	// Delete chirp
	if err := cfg.db.DeleteChirpByID(r.Context(), chirpID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
