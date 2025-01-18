package main

import (
	"encoding/json"
	"net/http"

	"github.com/amanRathod/chirpy/internal/database"
	"github.com/google/uuid"
)

type PolkaWebhooks struct {
	Event    string `json:"event"`
	Data struct {
		UserID    string `json:"user_id"`
	}
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		PolkaWebhooks
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, "no content")
		return
	}

	userIdString := params.Data.UserID
	userId, err := uuid.Parse(userIdString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	_, errs := cfg.db.UpdateUserToChirpyRed(r.Context(), database.UpdateUserToChirpyRedParams{
		IsChirpyRed: true,
		ID: userId,
	})
	if errs != nil {
		respondWithError(w, http.StatusNotFound, "not found", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "Updated")
}