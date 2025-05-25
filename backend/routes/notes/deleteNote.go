package routes

import (
	"backend/db"
	"backend/util"
	"encoding/json"
	"fmt"
	"net/http"
)

type DeleteNoteRequest struct {
	NoteID uint `json:"note_id"`
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	// Parse the request body
	var request DeleteNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the request body is not empty
	if !util.ValidateStruct(request) {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	noteRepo := db.NewBlockRepository(db.GetDB())

	// Attempt to delete the note
	if err := noteRepo.DeleteNoteBlocks(userID, request.NoteID); err != nil {
		if err.Error() == fmt.Sprintf("no blocks found for noteID %d and userID %d", request.NoteID, userID) {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting note", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Note deleted successfully"}`))
}
