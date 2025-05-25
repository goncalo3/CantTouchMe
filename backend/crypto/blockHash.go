package crypto

import (
	"backend/models"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
)

// BlockHash computes the hash of a block.
// Parameters:
// - block: the block to hash
// Returns: the Base64-encoded SHA-256 hash of the block, or an error if the block cannot be marshaled to JSON
func BlockHash(block models.Block) (string, error) {
	// Create a string representation of the block
	blockBytes, err := json.Marshal(block)
	if err != nil {
		return "", err
	}

	// Compute the SHA-256 hash of the block string
	hash := sha256.Sum256(blockBytes)

	// Convert the hash to a Base64 string
	hashBase64 := base64.StdEncoding.EncodeToString(hash[:])

	return hashBase64, nil
}

// VerifyBlockChain verifies the integrity of a blockchain.
// Parameters:
// - blocks: a slice of blocks representing the blockchain
// Returns: a boolean indicating whether the blockchain is valid, and an error if any block cannot be hashed
func VerifyBlockChain(blocks []models.Block) (bool, error) {
	if len(blocks) == 0 {
		return true, nil // An empty chain is considered valid
	}

	// Handle the first block
	if blocks[0].PrevHash != "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=" {
		log.Printf("Invalid first block: expected prev hash %s, got %s", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", blocks[0].PrevHash)
		return false, nil
	}

	for i := 1; i < len(blocks); i++ {
		prevBlock := blocks[i-1]
		currentBlock := blocks[i]

		// Calculate the hash of the previous block
		expectedPrevHash, err := BlockHash(prevBlock)
		if err != nil {
			log.Printf("Error hashing block %d: %v", i-1, err)
			return false, err
		}

		if currentBlock.PrevHash != expectedPrevHash {
			return false, nil // Invalid chain
		}
	}

	return true, nil // The chain is valid
}
