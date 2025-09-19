package main

import (
	"encoding/json"
	"net/http"

	"github.com/chumaachike/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID string `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid api key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "apikey does not match", err)
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid json", err)
		return
	}

	// Only handle "user.upgraded" events
	if params.Event == "user.upgraded" {
		userID, err := uuid.Parse(params.Data.UserID)
		if err == nil {
			_, err = cfg.db.UpgradeUser(r.Context(), userID)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "could not upgrade user", err)
				return
			}
		}
	}

	// Always return 204, regardless of event
	w.WriteHeader(http.StatusNoContent)
}
