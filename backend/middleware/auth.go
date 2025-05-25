package middleware

import (
	"backend/auth"
	"context"
	"encoding/json"
	"net/http"
)

// AuthMiddleware checks if the user is authenticated by verifying the JWT token from the HTTP-only cookie, sets the user ID in the context, and calls the next handler.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
			return
		}

		// Verify the token
		claims, err := auth.VerifyToken(cookie.Value)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid authentication token"})
			return
		}

		// Store user information in the context for later use
		userID := uint32(claims["user_id"].(float64))
		// passes the user ID to the next handler for later use
		ctx := context.WithValue(r.Context(), "UserID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
