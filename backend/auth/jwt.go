package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	JWTSecret     string
	JWTExpiration int // in seconds
)

// SetJWTConfig sets the JWT configuration values
// Parameters:
// - secret: the secret key used for signing JWT tokens
// - expiration: the expiration time for JWT tokens in seconds
func SetJWTConfig(secret string, expiration int) {
	JWTSecret = secret
	JWTExpiration = expiration
}

type JWTClaims struct {
	UserID uint32 `json:"user_id"`
	jwt.MapClaims
}

// Function that generates a JWT token from the user's ID and signs it with the user's secret key.
// Parameters:
// - user_id: the ID of the user for whom the token is being generated
// - sign_type: the preferred signing method ("hmac-sha256", "hmac-sha512")
// Returns: the signed JWT token as a string, or an error if the signing process fails
func GenerateJWTToken(user_id uint32, sign_type string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(JWTSecret)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Second * time.Duration(JWTExpiration))
	switch sign_type {
	case "hmac-sha256":
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(key)
	case "hmac-sha512":
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(key)
	default:
		return "", errors.New("the signing method is not valid for the JWT token signing")
	}
}

// Function that verifies a JWT token and returns its contents.
// Parameters:
// - incoming_token: the JWT token to be verified
// Returns: the claims contained in the token if valid, or an error if the token is invalid
func VerifyToken(incoming_token string) (jwt.MapClaims, error) {
	key, err := base64.StdEncoding.DecodeString(JWTSecret)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(incoming_token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid JWT token")
}
