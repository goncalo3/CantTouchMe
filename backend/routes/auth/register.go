package routes

import (
	"backend/db"
	"backend/models"
	"backend/util"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

// RegisterRequestBody represents the JSON body for registration
type RegisterRequestBody struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	HMACType       string `json:"hmac_type"`
	EncryptionType string `json:"encryption_type"`
	LoginSalt      string `json:"login_salt"`
	EncryptionSalt string `json:"encryption_salt"`
	HMACSalt       string `json:"hmac_salt"`
	PublicKey      string `json:"public_key"`
}

// RegisterResponseBody represents the JSON response for registration
type RegisterResponseBody struct {
	UserID  uint32 `json:"user_id"`
	Message string `json:"message"`
}

// RegisterHandler handles user registration via REST API
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Parse the request body
	var request RegisterRequestBody
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		writeJSONError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validateRegistrationRequest(request); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if email is already registered
	userRepo := db.NewUserRepository(db.GetDB())
	existingUser, err := userRepo.GetUserByEmail(strings.ToLower(request.Email))
	switch {
	case err != nil && err.Error() != "user not found":
		log.Printf("Error checking existing user: %v", err)
		writeJSONError(w, "Error checking existing user", http.StatusInternalServerError)
		return
	case existingUser != nil:
		writeJSONError(w, "Email already registered", http.StatusConflict)
		return
	}

	// Create user object
	user := models.User{
		Name:           request.Name,
		Email:          strings.ToLower(request.Email),
		PubKey:         request.PublicKey,
		LoginSalt:      request.LoginSalt,
		EncryptionSalt: request.EncryptionSalt,
		HMACSalt:       request.HMACSalt,
		HMACType:       request.HMACType,
		EncryptionType: request.EncryptionType,
	}

	// Save user to database
	userID, err := userRepo.CreateUser(&user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		if strings.Contains(err.Error(), "Duplicate entry") {
			writeJSONError(w, "Email already registered", http.StatusConflict)
			return
		}
		writeJSONError(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	response := RegisterResponseBody{
		UserID:  userID,
		Message: "User registered successfully",
	}

	// Write response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// Helper function to write JSON error responses
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// validateRegistrationRequest validates that all required fields are present and valid for registration
func validateRegistrationRequest(request RegisterRequestBody) error {
	if !util.ValidateStruct(request) {
		return errors.New("Required fields are missing")
	}

	// Check if all fields are valid
	switch {
	case !strings.Contains(request.Email, "@") || !strings.Contains(request.Email, "."):
		return errors.New("invalid email format")
	case request.HMACType != "hmac-sha256" && request.HMACType != "hmac-sha512":
		return errors.New("hmac type must be either hmac-sha256 or hmac-sha512")
	case request.EncryptionType != "aes-128-cbc" && request.EncryptionType != "aes-128-ctr":
		return errors.New("encryption type must be either aes-128-cbc or aes-128-ctr")
	case len(request.LoginSalt) < 44 || len(request.EncryptionSalt) < 44 || len(request.HMACSalt) < 44:
		return errors.New("salt must be encoded in base64")
	case len(request.PublicKey) < 44:
		return errors.New("public key must be encoded in base64")
	}

	return nil
}
