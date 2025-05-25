package routes

import (
	"backend/config"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConfig()

	// Determine SameSite policy based on environment
	sameSite := http.SameSiteStrictMode
	if cfg.Environment == "development" {
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: sameSite,
		MaxAge:   cfg.JWTExpiration,
	})

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
