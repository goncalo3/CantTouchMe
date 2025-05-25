package models

import (
	"time"
)

// This matches the structure of the block in frontend
type Block struct {
	PrevHash    string    `json:"prev_hash"`    // Hash of the previous block
	IV          string    `json:"iv"`           // Initialization vector for body encryption
	IVTitle     string    `json:"iv_title"`     // Initialization vector for title encryption
	CipherTitle string    `json:"cipher_title"` // Encrypted title
	Ciphertext  string    `json:"ciphertext"`   // Encrypted body content
	MAC         string    `json:"mac"`          // Message Authentication Code
	Signature   string    `json:"signature"`    // Digital signature of the block
	Timestamp   time.Time `json:"timestamp"`    // Block creation timestamp
}
