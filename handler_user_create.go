package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/amanRathod/chirpy/internal/database"
	"github.com/amanRathod/chirpy/internal/auth"
	"github.com/google/uuid"
)

type User struct {
	ID        	uuid.UUID `json:"id"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
	Email     	string    `json:"email"`
	IsChirpyRed bool    	`json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
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

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        		user.ID,
			CreatedAt: 		user.CreatedAt,
			UpdatedAt: 		user.UpdatedAt,
			Email:     		user.Email,
			IsChirpyRed: 	user.IsChirpyRed,
		},
	})
}

func (cfg *apiConfig) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		respondWithError(w, http.StatusForbidden, "You are not allowed to perform this action", nil)
		return
	}

	err := cfg.db.Reset(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, "Deleted all users")
}