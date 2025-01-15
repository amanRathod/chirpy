package main

import (
	"encoding/json" // Provides functionality for encoding and decoding JSON.
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	// Define a struct to parse the incoming JSON body.
	type parameters struct {
		Body string `json:"body"` // Map the JSON key "body" to the struct field `Body`.
	}

	// Define a struct to structure the response payload.
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	// Create a JSON decoder to read and parse the request body.
	decoder := json.NewDecoder(r.Body)

	// Initialize an instance of `parameters` to hold the parsed input.
	params := parameters{}

	// Decode the JSON body into the `params` struct.
	err := decoder.Decode(&params)
	if err != nil {
		// If decoding fails, respond with an error message and return.
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Define a constant for the maximum allowable chirp length.
	const maxChirpLength = 140

	// Check if the length of the "body" exceeds the maximum allowed chirp length.
	if len(params.Body) > maxChirpLength {
		// Respond with an error message if the chirp is too long.
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned,
	})
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}