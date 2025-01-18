package main

import (
	"encoding/json"
	"net/http"

	"github.com/amanRathod/chirpy/internal/auth"
	"github.com/amanRathod/chirpy/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	UserID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	email := params.Email
	password := params.Password

	if email == "" || password == "" {
		respondWithError(w, http.StatusBadRequest, "Please provide required data", nil)
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: email,
		HashedPassword: hashedPassword,
		ID: UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: 					user.ID,
			Email: 				user.Email,
			CreatedAt: 		user.CreatedAt,
			UpdatedAt: 		user.UpdatedAt,
			IsChirpyRed: 	user.IsChirpyRed,
		},
	})
}