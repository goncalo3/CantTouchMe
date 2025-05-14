package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"backend/config"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	
	// Define a simple Hello World handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		message := fmt.Sprintf("Hello World from %s environment!", cfg.Environment)
		w.Write([]byte(message))
	})
	
	// Start the server
	serverAddr := ":" + strconv.Itoa(cfg.Port)
	log.Printf("Server starting on port %d in %s mode with log level %s", 
		cfg.Port, cfg.Environment, cfg.LogLevel)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}