package crypto

import (
	"bytes"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/ed25519"

	"backend/models" // Corrected import path based on go.mod
)

// VerifyEd25519Signature verifies if a signature is valid for the given message and public key.
// Parameters:
// - publicKeyBase64: the Base64-encoded Ed25519 public key
// - messageBase64: the Base64-encoded message to verify
// - signatureBase64: the Base64-encoded signature to verify
// Returns: a boolean indicating whether the signature is valid, and an error if any input is invalid
func VerifyEd25519Signature(publicKeyBase64 string, messageBase64 string, signatureBase64 string) (bool, error) {
	// Decode the base64 strings
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return false, errors.New("invalid public key format")
	}

	messageBytes, err := base64.StdEncoding.DecodeString(messageBase64)
	if err != nil {
		return false, errors.New("invalid message format")
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false, errors.New("invalid signature format")
	}

	// Ed25519 public keys should be 32 bytes
	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return false, errors.New("invalid public key size")
	}

	// Ed25519 signatures should be 64 bytes
	if len(signatureBytes) != ed25519.SignatureSize {
		return false, errors.New("invalid signature size")
	}

	// Cast to the appropriate type
	publicKey := ed25519.PublicKey(publicKeyBytes)

	// Verify the signature using the ed25519 library
	return ed25519.Verify(publicKey, messageBytes, signatureBytes), nil
}

// verifyBlockEd25519Signature verifies the signature of a block.
// Parameters:
// - publicKeyBase64: the Base64-encoded Ed25519 public key
// - block: a pointer to the block whose signature is to be verified
// Returns: a boolean indicating whether the block's signature is valid, and an error if any input is invalid
func VerifyBlockEd25519Signature(publicKeyBase64 string, block *models.Block) (bool, error) {
	// Decode the public key
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return false, errors.New("invalid public key format")
	}

	// Concatenate the string
	// Prepare the data to verify
	var buffer bytes.Buffer
	string := block.PrevHash + block.IV + block.IVTitle + block.CipherTitle + block.Ciphertext + block.MAC + block.Timestamp.Format(time.RFC3339)

	log.Printf("Data to verify: %s", string)
	buffer.WriteString(string)

	dataToVerify := buffer.Bytes()

	// Decode the signature
	signatureBytes, err := base64.StdEncoding.DecodeString(block.Signature)
	if err != nil {
		return false, errors.New("invalid signature format")
	}

	// Verify the signature using the ed25519 library
	publicKey := ed25519.PublicKey(publicKeyBytes)
	return ed25519.Verify(publicKey, dataToVerify, signatureBytes), nil
}
