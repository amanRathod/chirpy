package main

import (
	"context"
	"net/http"

	"github.com/amanRathod/chirpy/internal/auth"
)



func (cfg *apiConfig) validateCredentialsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Please provide token", err)
			return
		}

		userID, errs := auth.ValidateJWT(token, cfg.jwtSecret)
		if errs != nil {
			respondWithError(w, http.StatusUnauthorized, "Token is invalid", err)
			return
		}

		// Add userID to the request context for downstream handlers
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", userID)
		r = r.WithContext(ctx)

		next(w,r)
	}
}