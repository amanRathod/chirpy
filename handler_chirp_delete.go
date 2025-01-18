package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	UserID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

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

	if dbChirp.UserID != UserID {
		respondWithError(w, http.StatusForbidden, "Not authorized", nil)
		return
	}

	errs := cfg.db.DeleteChirpById(r.Context(), dbChirp.ID)
	if errs != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong", errs)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "chirp deleted")
}