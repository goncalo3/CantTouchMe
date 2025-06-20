package routes

import (
	"backend/crypto"
	"backend/db"
	"backend/models"
	"backend/util"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Bassically the same logic as add block but with an id generated by the db
type NewNoteResponse struct {
	TimeStamp string `json:"timestamp"`
	NoteID    uint   `json:"note_id"`
	Message   string `json:"message"`
}

func NewNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Decode the request payload into NewNoteRequest
	var request models.Block
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request body
	if !util.ValidateStruct(request) {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	// Check if the hash is properly initialized
	if request.PrevHash != "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=" {
		http.Error(w, "Invalid request body: PrevHash not properly be initialized", http.StatusBadRequest)
		return
	}

	// set the user ID usig the jwt middleware
	userID, ok := r.Context().Value("UserID").(uint32)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
	isValid, err := crypto.VerifyBlockEd25519Signature(user.PubKey, &request)
	if err != nil || !isValid {
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	blockRepo := db.NewBlockRepository(db.GetDB())

	// Create a new note in the database
	NoteId, err := blockRepo.CreateNewNote(userID, &request)
	if err != nil {
		log.Printf("Error creating new note: %v", err)
		http.Error(w, "Error creating block", http.StatusInternalServerError)
		return
	}

	response := NewNoteResponse{
		TimeStamp: request.Timestamp.Format(time.RFC3339),
		NoteID:    NoteId,
		Message:   "Block created successfully!",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
