package routes

import (
	"backend/auth"
	"backend/config"
	"backend/crypto"
	"backend/db"
	"backend/models"
	"backend/util"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// LoginRequestBody represents the JSON body of a login request
type LoginRequestBody struct {
	// in the unlikely event of collisions, we also have the email
	Email     string `json:"email"`
	Challenge string `json:"challenge"`
	Signature string `json:"signature"`
}

// LoginResponseBody represents the JSON response for a successful login
// Returns the JWT token, user ID, the HMAC salt, and the encryption salt inside the user object
type LoginResponseBody struct {
	Message string      `json:"message"`
	User    models.User `json:"user"`
}

// LoginHandler verifies the signed challenge and issues a JWT token on success
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var requestBody LoginRequestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// valide all fields are present
	if !util.ValidateStruct(requestBody) {
		http.Error(w, "Challenge, signature  and Email are required", http.StatusBadRequest)
		return
	}

	// Validate the request body fields
	switch {
	case !strings.Contains(requestBody.Email, "@") || !strings.Contains(requestBody.Email, "."):
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	case len(requestBody.Signature) != 88:
		http.Error(w, "Invalid signature format", http.StatusBadRequest)
		return
	case len(requestBody.Challenge) != 44:
		http.Error(w, "Invalid challenge format", http.StatusBadRequest)
		return
	}

	// Get the challenge from the database for the given email and challenge
	challengeRepo := db.NewChallengeRepository(db.GetDB())
	challenge, err := challengeRepo.GetChallenge(requestBody.Challenge, requestBody.Email)
	if err != nil {
		log.Printf("Challenge lookup error: %v", err)
		http.Error(w, "Challenge not found", http.StatusUnauthorized)
		return
	}

	// Check if challenge is expired or has already been used
	switch {
	case time.Now().After(challenge.ExpiresAt):
		http.Error(w, "Challenge expired", http.StatusUnauthorized)
		return
	case challenge.Used:
		http.Error(w, "Challenge already used", http.StatusUnauthorized)
		return
	}

	// Get the user associated with the challenge
	userRepo := db.NewUserRepository(db.GetDB())
	user, err := userRepo.GetUserByID(challenge.UserID)
	if err != nil {
		log.Printf("User lookup error: %v", err)
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Verify the signature against the user's public key
	valid, err := crypto.VerifyEd25519Signature(
		user.PubKey,
		requestBody.Challenge, // We pass the challenge directly as its already be base64 encoded
		requestBody.Signature,
	)

	if err != nil {
		log.Printf("Signature verification error: %v", err)
		http.Error(w, "Signature verification failed", http.StatusUnauthorized)
		return
	}

	if !valid {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Delete the used challenge
	err = challengeRepo.DeleteChallenge(challenge.ID)
	if err != nil {
		log.Printf("Failed to delete used challenge: %v", err)
		http.Error(w, "Failed to delete used challenge", http.StatusInternalServerError)
		return
	}

	// Generate a JWT token
	token, err := auth.GenerateJWTToken(user.ID, user.HMACType)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	config := config.LoadConfig()

	// Determine SameSite policy based on environment
	sameSite := http.SameSiteStrictMode
	if config.Environment == "development" {
		sameSite = http.SameSiteNoneMode
	}

	// Set the JWT token in an HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: sameSite,
		MaxAge:   config.JWTExpiration,
	})

	// Return success response with the user data
	response := LoginResponseBody{
		Message: "Login successful",
		User:    *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
