package cron

import (
	"backend/config"
	"backend/db"
	"log"
	"time"
)

// CronScheduler manages scheduled tasks
type CronScheduler struct {
	running bool
	stopCh  chan bool
}

// NewCronScheduler creates a new cron scheduler
func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		running: false,
		stopCh:  make(chan bool),
	}
}

// Start begins the cron scheduler
func (cs *CronScheduler) Start() {
	if cs.running {
		log.Println("Cron scheduler is already running")
		return
	}

	cs.running = true
	log.Println("Starting cron scheduler...")

	go cs.run()
}

// Stop stops the cron scheduler
func (cs *CronScheduler) Stop() {
	if !cs.running {
		return
	}

	log.Println("Stopping cron scheduler...")
	cs.stopCh <- true
	cs.running = false
}

// run contains the main cron loop
func (cs *CronScheduler) run() {
	// Run cleanup immediately on startup
	cs.cleanupExpiredChallenges()

	// Get cleanup interval from config
	cfg := config.GetConfig()
	cleanupInterval := time.Duration(cfg.ChallengeCleanupMinutes) * time.Minute

	// Create a ticker with the configured interval
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	log.Printf("Challenge cleanup cron job scheduled to run every %d minutes", cfg.ChallengeCleanupMinutes)

	for {
		select {
		case <-ticker.C:
			cs.cleanupExpiredChallenges()
		case <-cs.stopCh:
			log.Println("Cron scheduler stopped")
			return
		}
	}
}

// cleanupExpiredChallenges removes expired challenges from the database
func (cs *CronScheduler) cleanupExpiredChallenges() {
	log.Println("Running scheduled cleanup of expired challenges...")

	challengeRepo := db.NewChallengeRepository(db.GetDB())
	err := challengeRepo.ClearExpiredChallenges()
	if err != nil {
		log.Printf("Error cleaning up expired challenges: %v", err)
		return
	}

	log.Println("Successfully cleaned up expired challenges")
}
