package models

import (
	"time"
)

// type note is note_id and block
type Note struct {
	NoteID uint  `json:"note_id"` // Unique identifier for the note
	Block  Block `json:"block"`   // The block data associated with the note
}

type NoteBlockChain struct {
	NoteID uint    `json:"note_id"` // Unique identifier for the note
	Blocks []Block `json:"blocks"`  // List of blocks in the note's blockchain
}

// For the endpoint that returns the titles
type Title struct {
	NoteID      uint      `json:"note_id"`
	CipherTitle string    `json:"cipher_title"`
	Timestamp   time.Time `json:"timestamp"`
	IV          string    `json:"iv_title"`
}
