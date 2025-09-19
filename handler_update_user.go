package main

import (
	"encoding/json"
	"net/http"

	"github.com/chumaachike/chirpy/internal/auth"
	"github.com/chumaachike/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    *string `json:"email"`    // pointers let you detect missing fields
		Password *string `json:"password"` // same here
	}
	type response struct {
		User
	}

	// Decode body
	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid json", err)
		return
	}

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

	// Hash password if provided
	var hashedPassword string
	if params.Password != nil {
		hashedPassword, err = auth.HashPassword(*params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "unable to hash password", err)
			return
		}
	}

	// Run update
	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          derefString(params.Email),
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to update user", err)
		return
	}

	// Send response
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			Email:       user.Email,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
