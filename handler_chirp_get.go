package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/amanRathod/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sorting := r.URL.Query().Get("sort")

	chirps := []Chirp{}
	var err error

	if authorID != "" {
		chirps, err = cfg.getChirpsByAuthor(r.Context(), authorID)
		} else {
				chirps, err = cfg.getAllChirps(r.Context())
		}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps", err)
		return
	}

	if sorting == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpsByAuthor(ctx context.Context, authorID string) ([]Chirp, error) {
	parsedAuthorID, err := parseUUID(authorID)
	if err != nil {
			return nil, fmt.Errorf("invalid author ID: %w", err)
	}

	dbChirps, err := cfg.db.GetChirpByUserId(ctx, parsedAuthorID)
	if err != nil {
			return nil, fmt.Errorf("database error retrieving chirps by user: %w", err)
	}

	return transformChirps(dbChirps), nil
}

func (cfg *apiConfig) getAllChirps(ctx context.Context) ([]Chirp, error) {
	dbChirps, err := cfg.db.GetChirps(ctx)
	if err != nil {
			return nil, fmt.Errorf("database error retrieving all chirps: %w", err)
	}

	return transformChirps(dbChirps), nil
}

func (cfg *apiConfig) handlerChirpById(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	})
}

func transformChirps(dbChirps []database.Chirp) []Chirp {
	chirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
			chirps[i] = Chirp{
					ID:        dbChirp.ID,
					CreatedAt: dbChirp.CreatedAt,
					UpdatedAt: dbChirp.UpdatedAt,
					UserID:    dbChirp.UserID,
					Body:      dbChirp.Body,
			}
	}
	return chirps
}


func parseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}