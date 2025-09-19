package main

import (
	"net/http"

	"github.com/chumaachike/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	var dbChirps []database.Chirp
	var err error

	if authorIDStr != "" {
		authorID, parseErr := uuid.Parse(authorIDStr)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "invalid author_id", parseErr)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByAuthor(r.Context(), database.GetChirpsByAuthorParams{
			UserID:    authorID,
			SortOrder: normalizeSort(sortOrder),
		})
	} else {
		dbChirps, err = cfg.db.GetAllChirps(r.Context(), normalizeSort(sortOrder))
	}
	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not fetch chirps", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func normalizeSort(s string) string {
	if s == "desc" {
		return "desc"
	}
	return "asc"
}
