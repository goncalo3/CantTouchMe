package routes

import (
	"backend/db"
	"backend/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

// UpdateUserRequestBody represents the JSON body for user updates
type UpdateUserRequestBody struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// UpdateUserResponseBody represents the JSON response for user updates
type UpdateUserResponseBody struct {
	Message string      `json:"message"`
	User    models.User `json:"user"`
}

// UpdateUserHandler handles user information updates
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept PUT requests
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("UserID").(uint32)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body
	var request UpdateUserRequestBody
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		writeJSONError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validateUpdateRequest(request); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user repository
	userRepo := db.NewUserRepository(db.GetDB())

	// Get current user
	currentUser, err := userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		writeJSONError(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	// Update user fields if provided
	if request.Name != "" {
		currentUser.Name = request.Name
	}
	if request.Email != "" {
		// Check if email is already taken by another user
		existingUser, err := userRepo.GetUserByEmail(strings.ToLower(request.Email))
		if err == nil && existingUser.ID != userID {
			writeJSONError(w, "Email already in use", http.StatusConflict)
			return
		}
		currentUser.Email = strings.ToLower(request.Email)
	}

	// Update user in database
	err = userRepo.UpdateUser(currentUser)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		writeJSONError(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := UpdateUserResponseBody{
		Message: "User updated successfully",
		User:    *currentUser,
	}

	// Write response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// validateUpdateRequest validates the update request fields
func validateUpdateRequest(request UpdateUserRequestBody) error {
	if request.Email != "" {
		if !strings.Contains(request.Email, "@") || !strings.Contains(request.Email, ".") {
			return errors.New("invalid email format")
		}
	}
	return nil
}
