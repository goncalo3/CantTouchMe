package middleware

import (
	"backend/config" // Import the config package
	"net/http"
	"strings"
)

// CorsMiddleware adds CORS headers to HTTP responses to allow cross-origin requests.
func CorsMiddleware(next http.Handler) http.Handler {
	cfg := config.GetConfig() // Retrieve the application configuration

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin") // Get the Origin header from the request

		allowedOrigin := "" // Initialize the allowed origin as an empty string
		if cfg.Environment == "production" {
			// In production, only allow requests from the specified domain
			allowedOrigin = "https://cantotuchme.goncalo3.pt"
		} else {
			// In non-production environments, allow requests from localhost
			if origin != "" && (origin == "http://localhost" ||
				strings.HasPrefix(origin, "http://localhost:")) {
				allowedOrigin = origin
			}
		}

		if allowedOrigin != "" {
			// Set CORS headers if the origin is allowed
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Set headers for CORS preflight requests
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Set allowed headers for CORS requests
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			// Handle preflight requests by returning a 200 OK status
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass the request to the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
