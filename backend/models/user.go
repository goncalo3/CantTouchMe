package models

// represents a user in the system.
type User struct {
	ID             uint32 `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	PubKey         string `json:"public_key"`
	EncryptionSalt string `json:"encryption_salt"`
	HMACSalt       string `json:"hmac_salt"`
	HMACType       string `json:"hmac_type"`
	EncryptionType string `json:"encryption_type"`
	LoginSalt      string `json:"login_salt"`
}
