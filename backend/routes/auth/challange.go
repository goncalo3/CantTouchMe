package routes

import (
	"backend/crypto"
	"backend/db"
	"backend/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// ChallengeRequestBody represents the JSON body of a challenge request
type ChallengeRequestBody struct {
	Email string `json:"email"`
}

// ChallengeResponseBody represents the JSON response with the challenge
type ChallengeResponseBody struct {
	Challenge string `json:"challenge"`
	LoginSalt string `json:"login_salt"`
	ExpiresAt string `json:"expires_at"`
}

// ChallengeHandler generates and sends an authentication challenge
func ChallengeHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var requestBody ChallengeRequestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate email
	if requestBody.Email == "" || (!strings.Contains(requestBody.Email, "@") && !strings.Contains(requestBody.Email, ".")) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	// Generate a random challenge value
	challengeValue, err := crypto.GenerateSaltBase64(32)
	if err != nil {
		log.Printf("Error generating challenge: %v", err)
		http.Error(w, "Failed to generate challenge", http.StatusInternalServerError)
		return
	}

	// Set challenge expiration (5 minutes from now)
	expiresAt := time.Now().Add(5 * time.Minute)

	// Get user by email
	userRepo := db.NewUserRepository(db.GetDB())
	user, err := userRepo.GetUserByEmail(strings.ToLower(requestBody.Email))
	if err != nil {
		// even if no user is found, we still create a challenge with a random salt and send it, this is to prevent user enumeration attacks
		dummySalt, err := crypto.GenerateSaltBase64(32)
		if err != nil {
			log.Printf("Error generating dummy salt: %v", err)
			http.Error(w, "Unknown error ocurred when generating challange", http.StatusInternalServerError)
			return
		}
		response := ChallengeResponseBody{
			Challenge: string(challengeValue),
			LoginSalt: dummySalt,
			ExpiresAt: expiresAt.Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create challenge in the database
	challenge := models.Challenge{
		UserID:         user.ID,
		ChallengeValue: string(challengeValue),
		ExpiresAt:      expiresAt,
		Used:           false,
	}

	// Save challenge
	challengeRepo := db.NewChallengeRepository(db.GetDB())
	_, err = challengeRepo.CreateChallenge(&challenge)
	if err != nil {
		log.Printf("Error creating challenge: %v", err)
		http.Error(w, "Failed to create challenge", http.StatusInternalServerError)
		return
	}

	// Respond with the challenge
	response := ChallengeResponseBody{
		Challenge: string(challengeValue),
		LoginSalt: user.LoginSalt,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
