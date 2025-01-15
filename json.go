package main

import (
	"encoding/json" // Provides functionality for encoding and decoding JSON.
	"log"
	"net/http"
)

// respondWithError is a utility function to send error responses in JSON format.
// It logs the error if present and sends a JSON-formatted error message as the response.
func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	// Log the error details if an error object is provided.
	if err != nil {
		log.Println(err) // Print the error message to the server logs.
	}

	// Log 5XX server errors for further analysis.
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	// Define a struct to format the error message in JSON response.
	type errorResponse struct {
		Error string `json:"error"` // Map the JSON key "error" to the struct field `Error`.
	}

	// Send the error response as JSON with the given HTTP status code.
	respondWithJSON(w, code, errorResponse{
		Error: msg, // Pass the error message to the response struct.
	})
}

// respondWithJSON is a utility function to send any response in JSON format.
// It sets the Content-Type header, marshals the payload into JSON, and writes it to the response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Set the response Content-Type header to indicate JSON data.
	w.Header().Set("Content-Type", "application/json")

	// Marshal the payload into a JSON byte slice.
	dat, err := json.Marshal(payload)
	if err != nil {
		// Log the error if JSON marshaling fails.
		log.Printf("Error marshalling JSON: %s", err)

		// Respond with a 500 Internal Server Error if marshaling fails.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the HTTP status code to the response header.
	w.WriteHeader(code)

	// Write the JSON-encoded payload to the response body.
	w.Write(dat)
}
