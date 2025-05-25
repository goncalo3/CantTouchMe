package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// GenerateSalt generates a cryptographically secure random salt.
// Parameters:
// - size: the size of the salt in bytes
// Returns: a byte slice containing the generated salt, or an error if the salt generation fails
func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, salt)
	return salt, err
}

// GenerateSaltBase64 generates a cryptographically secure random salt and encodes it in Base64.
// Parameters:
// - size: the size of the salt in bytes
// Returns: a Base64-encoded string of the generated salt, or an error if the salt generation fails
func GenerateSaltBase64(size int) (string, error) {
	salt, err := GenerateSalt(size)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}
