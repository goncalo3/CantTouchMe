package models

import "time"

// Challenge represents a login challenge for public key authentication.
type Challenge struct {
	ID             uint32    `json:"id"`
	UserID         uint32    `json:"user_id"`
	ChallengeValue string    `json:"challenge_value"`
	ExpiresAt      time.Time `json:"expires_at"`
	Used           bool      `json:"used"`
	CreatedAt      time.Time `json:"created_at"`
}
