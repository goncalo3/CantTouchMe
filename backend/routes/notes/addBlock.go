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

// AddBlockRquest represents the request body for adding a block

type AddBlockResponse struct {
	TimeStamp string `json:"timestamp"`
	Message   string `json:"message"`
}

// this can also be tought of as the edit note endpoint
func AddBlockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// The request format is the block
	var request models.Note
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request body
	if util.ValidateStruct(request) == false {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	// get the user ID from the context set by the JWT middleware
	userID, ok :=
		r.Context().Value("UserID").(uint32)
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
	isValid, err := crypto.VerifyBlockEd25519Signature(user.PubKey, &request.Block)
	if err != nil || !isValid {
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	blockRepo := db.NewBlockRepository(db.GetDB())

	// Get all blocks to determine the previous hash
	blockchain, err := blockRepo.GetNoteBlockChain(userID, request.NoteID)
	if err != nil {
		log.Printf("Error retrieving blocks for user %d and note %d: %v", userID, request.NoteID, err)
		http.Error(w, "Error retrieving blocks", http.StatusInternalServerError)
		return
	}

	// Checks the blockhain integrity
	// If the blockchain is invalid, it will become impossible to edit the note
	valid, err := crypto.VerifyBlockChain(blockchain.Blocks)
	if err != nil || !valid {
		log.Printf("Invalid block chhain for note %d: %v! You can no longer edit this note!", request.NoteID, err)
		http.Error(w, "Invalid block chain! You can no longer edit this note!", http.StatusBadRequest)
		return
	}

	// Create the block in the database
	err = blockRepo.CreateBlock(userID, request.NoteID, &request.Block)
	if err != nil {
		log.Printf("Error creating block for user %d and note %d: %v", userID, request.NoteID, err)
		http.Error(w, "Error creating block", http.StatusInternalServerError)
		return
	}

	response := AddBlockResponse{
		TimeStamp: request.Block.Timestamp.Format(time.RFC3339),
		Message:   "Block created successfully!",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
