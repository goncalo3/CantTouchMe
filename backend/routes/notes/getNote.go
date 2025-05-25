package routes

import (
	"backend/crypto"
	"backend/db"
	"backend/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// GetNotesRequest defines the JSON shape for the note request from clients
type GetNotesRequest struct {
	NoteID uint `json:"note_id"`
}

// Returns the latest block of a note
func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	// The request format is the block
	var request GetNotesRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !util.ValidateStruct(request) {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	blockRepo := db.NewBlockRepository(db.GetDB())

	// Fetch the note block for the given userID and noteID
	block, err := blockRepo.GetNoteBlock(userID, request.NoteID)
	if err != nil {
		if err.Error() == fmt.Sprintf("no block found for noteID %d and userID %d", request.NoteID, userID) {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			log.Printf("Error retrieving note: %v", err)
			http.Error(w, "Error retrieving note", http.StatusInternalServerError)
		}
		return
	}

	// Fetch the public key from the database using the user ID
	userRepo := db.NewUserRepository(db.GetDB())
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Check if the signature is valid
	isValid, err := crypto.VerifyBlockEd25519Signature(user.PubKey, block)
	if err != nil || !isValid {
		http.Error(w, "Invalid signature!", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(block)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
