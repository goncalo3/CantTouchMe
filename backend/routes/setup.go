package routes

import (
	"backend/middleware"
	auth "backend/routes/auth"
	notes "backend/routes/notes"
	"net/http"
)

// SetupAuthRoutes registers all the authentication routes with the router
func SetupAuthRoutes(mux *http.ServeMux) {
	// Register route
	mux.HandleFunc("/auth/register", auth.RegisterHandler)

	// Challenge route
	mux.HandleFunc("/auth/challenge", auth.ChallengeHandler)

	// Login route
	mux.HandleFunc("/auth/login", auth.LoginHandler)

	// Logout route
	mux.HandleFunc("/auth/logout", auth.LogoutHandler)

	// Delete user route
	mux.HandleFunc("/auth/delete", middleware.AuthMiddleware(auth.DeleteUserHandler))
	// Update user route
	mux.HandleFunc("/auth/update", middleware.AuthMiddleware(auth.UpdateUserHandler))

	// Note edition adds a new block to the note blockchain
	mux.HandleFunc("/notes/edit", middleware.AuthMiddleware(notes.AddBlockHandler))

	// Creates a new note from 0
	mux.HandleFunc("/notes/new", middleware.AuthMiddleware(notes.NewNoteHandler))

	// get all the notes titles
	mux.HandleFunc("/notes/titles", middleware.AuthMiddleware(notes.GetTitlesHandler))

	// get note by id
	mux.HandleFunc("/notes/get", middleware.AuthMiddleware(notes.GetNoteHandler))

	// delete a note by id
	mux.HandleFunc("/notes/delete", middleware.AuthMiddleware(notes.DeleteNoteHandler))
}
