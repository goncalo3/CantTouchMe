package routes

import (
	"backend/db"
	"backend/models"
	"encoding/json"
	"log"
	"net/http"
)

// GetTitlesHandler gets titles of all the notes for a user
func GetTitlesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Extract userID from context
	userID, ok := r.Context().Value("UserID").(uint32)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	blockRepo := db.NewBlockRepository(db.GetDB())

	// Fetch titles for the user
	titles, err := blockRepo.GetTitles(userID)
	if err != nil {
		log.Printf("Error retrieving titles: %v", err)
		http.Error(w, "Error retrieving titles", http.StatusInternalServerError)
		return
	}

	// Check if titles is nil and send an empty array if so
	if titles == nil {
		titles = []*models.Title{}
	}

	err = json.NewEncoder(w).Encode(titles)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
