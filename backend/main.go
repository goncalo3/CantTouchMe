package main

import (
	"backend/auth"
	"backend/config"
	"backend/cron"
	"backend/db"
	"backend/middleware"
	routes "backend/routes"
	"log"
	"net/http"
	"strconv"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	dbCfg := config.LoadDbConfig()

	// Set JWT configuration
	auth.SetJWTConfig(cfg.JWTSecret, cfg.JWTExpiration)

	// Initialize database connection
	db.InitDB(dbCfg)
	defer db.CloseDB()

	// Start cron scheduler for cleanup tasks
	cronScheduler := cron.NewCronScheduler()
	cronScheduler.Start()
	defer cronScheduler.Stop()

	// Create a new ServeMux for routing
	mux := http.NewServeMux()

	// Set up auth routes
	routes.SetupAuthRoutes(mux)

	handler := middleware.CorsMiddleware(mux)

	// Apply logging middleware only in development environment
	if cfg.Environment == "development" {
		handler = middleware.LogMiddleware()(middleware.CorsMiddleware(mux))
	}

	// Start the server
	serverAddr := ":" + strconv.Itoa(cfg.Port)
	log.Printf("CantTouchMe server starting on port %d in %s mode!",
		cfg.Port, cfg.Environment)
	err := http.ListenAndServe(serverAddr, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
